package main

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/repository"
	"github.com/ArtyomArtamonov/msg/internal/server"
	"github.com/ArtyomArtamonov/msg/internal/service"
	"github.com/streadway/amqp"

	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.TraceLevel)

	err := godotenv.Load("../../.env")
	failOnError(err, "Error loading .env file")

	env := server.NewEnv()

	host := env.API_HOST
	lis, err := net.Listen("tcp", host)
	failOnError(err, "could not create tcp connection")

	connectionString := fmt.Sprintf(
		"host=database port=5432 sslmode=disable dbname=%s user=%s password=%s",
		env.POSTGRES_DB,
		env.POSTGRES_USER,
		env.POSTGRES_PASSWORD,
	)
	db, err := sql.Open("postgres", connectionString)
	failOnError(err, "could not connect to database")

	err = db.Ping()
	failOnError(err, "could not ping database")
	defer db.Close()

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@message-broker:5672/",
		env.RABBITMQ_DEFAULT_USER,
		env.PGADMIN_DEFAULT_PASSWORD))
	failOnError(err, "could not connect to message-broker")

	ch, err := conn.Channel()
	failOnError(err, "could not create channel (message-broker)")
	defer ch.Close()

	grpcServer := createAndPrepareGRPCServer(db, ch, env)

	logrus.Info("Starting grpc server on ", host)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal(err.Error())
	}
}

func createAndPrepareGRPCServer(db *sql.DB, ch *amqp.Channel, env *server.Env) *grpc.Server {
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
		if err := userStore.Save(user); err != nil {
			logrus.Error(err)
		}

		admin, err := model.NewUser("admin", "admin", model.ADMIN_ROLE)
		if err != nil {
			logrus.Error(err)
		}
		if err := userStore.Save(admin); err != nil {
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

	// API
	amqpManager := service.NewRabbitMQManager(ch)
	roomStore := repository.NewPostgresRoomStore(db)
	apiServer := server.NewApiServer(jwtManager, roomStore, amqpManager)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterApiServiceServer(grpcServer, apiServer)

	reflection.Register(grpcServer)

	return grpcServer
}

func failOnError(err error, text string) {
	if err != nil {
		logrus.Fatalf("%s: %v", text, err)
	}
}
