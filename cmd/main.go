package main

import (
	"net"
	"os"
	"strconv"
	"time"

	authPb "github.com/ArtyomArtamonov/msg/pkg/auth/proto"
	messagePb "github.com/ArtyomArtamonov/msg/pkg/message/proto"

	"github.com/ArtyomArtamonov/msg/pkg/auth"
	"github.com/ArtyomArtamonov/msg/pkg/message"
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

	grpcServer := createAndPrepareGRPCServer()

	logrus.Info("Starting grpc server on ", host)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal(err.Error())
	}
}

func createAndPrepareGRPCServer() *grpc.Server {
	jwtDurationMin, err := strconv.Atoi(os.Getenv("JWT_DURATION_MIN"))
	if err != nil {
		logrus.Fatal("Could not get JWT_DURATION_MIN env variable (should be a number of minutes token expiration time)")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	userStore := auth.NewInMemoryUserStore()
	// DEBUG PURPOSE BLOCK
	{
		user, err := auth.NewUser("user", "user", "user")
		if err != nil {
			logrus.Fatal("Could not create test user")
		}
		userStore.Save(user)

		admin, err := auth.NewUser("admin", "admin", "admin")
		if err != nil {
			logrus.Fatal("Could not create test admin")
		}
		userStore.Save(admin)
	}
	jwtManager := auth.NewJWTManager(jwtSecret, time.Minute*time.Duration(jwtDurationMin))

	authServer := auth.NewAuthService(userStore, jwtManager)
	messageServer := message.NewMessageService(jwtManager)

	authInterceptor := auth.NewAuthInterceptor(jwtManager, accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	messagePb.RegisterMessageServiceServer(grpcServer, messageServer)
	authPb.RegisterAuthServiceServer(grpcServer, authServer)

	reflection.Register(grpcServer)

	return grpcServer
}

func accessibleRoles() map[string][]string {
	const authService = "/message.MessageService/"
	const messageService = "/message.MessageService/"

	return map[string][]string{
		messageService + "SendMessage": {"admin", "user"},
		messageService + "GetMessages": {"admin", "user"},
	}
}
