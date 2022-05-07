package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/repository"
	"github.com/ArtyomArtamonov/msg/internal/server"
	"github.com/ArtyomArtamonov/msg/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgtypeuuid "github.com/vgarvardt/pgx-google-uuid/v5"

	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.TraceLevel)

	err := godotenv.Load("../.env")
	if err != nil {
		logrus.Fatal("Error loading .env file: ", err)
	}
	env := server.NewEnv()

	host := env.HOST
	lis, err := net.Listen("tcp", host)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@database:5432/%s",

		env.POSTGRES_USER,
		env.POSTGRES_PASSWORD,
		env.POSTGRES_DB,
	)

	dbconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		logrus.Fatal("could not parse database config")
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgtypeuuid.Register(conn.TypeMap())
		return nil
	}
	db, err := pgxpool.ConnectConfig(context.Background(), dbconfig)
	if err != nil {
		logrus.Fatal("could not connect to database")
	}
	if err := db.Ping(context.Background()); err != nil {
		logrus.Fatalf("could not ping database: %v", err)
	}
	defer db.Close()

	grpcServer := createAndPrepareGRPCServer(db, env)

	logrus.Info("Starting grpc server on ", host)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal(err.Error())
	}
}

func createAndPrepareGRPCServer(db *pgxpool.Pool, env *server.Env) *grpc.Server {
	endpoints := server.NewEndpoints()
	endpointRoles := server.NewEndpointRoles(endpoints)

	// AUTH
	userStore := repository.NewPostgresUserStore(db)
	refreshTokenStore := repository.NewRefreshTokenPostgresStore(db)
	// DEBUG PURPOSE BLOCK
	{
		user, err := model.NewUser("user", "user", model.USER_ROLE)
		if err != nil {
			logrus.Error(err)
		}
		//SHOULD BE REMOVED IN FUTURE
		if err := userStore.Save(context.Background(), user); err != nil {
			logrus.Error(err)
		}

		admin, err := model.NewUser("admin", "admin", model.ADMIN_ROLE)
		if err != nil {
			logrus.Error(err)
		}
		//SHOULD BE REMOVED IN FUTURE
		if err := userStore.Save(context.Background(), admin); err != nil {
			logrus.Error(err)
		}
	}
	jwtManager := service.NewJWTManager(
		env.JWT_SECRET,
		time.Minute*time.Duration(env.JWT_DURATION_MIN),
		time.Hour*24*time.Duration(env.REFRESH_DURATION_DAYS),
	)

	authServer := server.NewAuthServer(userStore, refreshTokenStore, jwtManager)
	authInterceptor := server.NewAuthInterceptor(jwtManager, endpointRoles)

	// MESSAGE
	sessionStore := repository.NewInMemorySessionStore()
	messageServer := server.NewMessageServer(jwtManager, sessionStore)

	// API
	roomStore := repository.NewPostgresRoomStore(db)
	apiServer := server.NewApiServer(jwtManager, roomStore)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	pb.RegisterMessageServiceServer(grpcServer, messageServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterApiServiceServer(grpcServer, apiServer)

	reflection.Register(grpcServer)

	return grpcServer
}
