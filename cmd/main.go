package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/repository"
	"github.com/ArtyomArtamonov/msg/internal/server"
	"github.com/ArtyomArtamonov/msg/internal/service"

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

	host := os.Getenv("HOST")
	lis, err := net.Listen("tcp", host)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	connectionString := fmt.Sprintf(
		"host=database port=5432 sslmode=disable dbname=%s user=%s password=%s",
		os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"))
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logrus.Fatal("could not connect to database")
	}
	if err := db.Ping(); err != nil {
		logrus.Fatalf("could not ping database: %v", err)
	}
	defer db.Close()

	grpcServer := createAndPrepareGRPCServer(db)

	logrus.Info("Starting grpc server on ", host)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal(err.Error())
	}
}

func createAndPrepareGRPCServer(db *sql.DB) *grpc.Server {
	jwtDurationMin, err := strconv.Atoi(os.Getenv("JWT_DURATION_MIN"))
	if err != nil {
		logrus.Fatal("Could not get JWT_DURATION_MIN env variable (should be a number of minutes token expiration time)")
	}

	refreshDurationDays, err := strconv.Atoi(os.Getenv("REFRESH_DURATION_DAYS"))
	if err != nil {
		logrus.Fatal("Could not get REFRESH_DURATION_DAYS env variable (should be a number of days refresh token expiration time)")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

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
	jwtManager := service.NewJWTManager(jwtSecret, time.Minute*time.Duration(jwtDurationMin), time.Hour*24*time.Duration(refreshDurationDays))
	
	authServer := server.NewAuthServer(userStore, refreshTokenStore, jwtManager)
	authInterceptor := server.NewAuthInterceptor(jwtManager, accessibleRoles())

	// MESSAGE
	sessionStore := repository.NewInMemorySessionStore()
	messageServer := server.NewMessageServer(jwtManager, sessionStore)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	pb.RegisterMessageServiceServer(grpcServer, messageServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	reflection.Register(grpcServer)

	return grpcServer
}

func accessibleRoles() map[string][]string {
	const authService = "/message.MessageService/"
	const messageService = "/message.MessageService/"

	return map[string][]string{
		messageService + "SendMessage": {model.ADMIN_ROLE, model.USER_ROLE},
		messageService + "GetMessages": {model.ADMIN_ROLE, model.USER_ROLE},
	}
}
