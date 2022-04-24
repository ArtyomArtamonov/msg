package main

import (
	"net"
	"os"
	"strconv"
	"time"

	apiPb "github.com/ArtyomArtamonov/msg/pkg/api/proto"
	authPb "github.com/ArtyomArtamonov/msg/pkg/auth/proto"
	"github.com/joho/godotenv"

	"github.com/ArtyomArtamonov/msg/pkg/api"
	"github.com/ArtyomArtamonov/msg/pkg/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const HOST = ":8000"

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		logrus.Fatal("Error loading .env file: ", err)
	}
	lis, err := net.Listen("tcp", HOST)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	grpcServer := createAndPrepareGRPCServer()

	logrus.Info("Starting grpc server on ", HOST)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal(err.Error())
	}
}

func createAndPrepareGRPCServer() *grpc.Server {
	apiServer := api.NewMessageService()

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

	authInterceptor := auth.NewAuthInterceptor(jwtManager, accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	apiPb.RegisterMessageServiceServer(grpcServer, apiServer)
	authPb.RegisterAuthServiceServer(grpcServer, authServer)

	reflection.Register(grpcServer)

	return grpcServer
}

func accessibleRoles() map[string][]string {
	const authService = "/api.MessageService/"
	const messageService = "/api.MessageService/"

	return map[string][]string{
		messageService + "SendMessage": {"admin", "user"},
		messageService + "GetMessages": {"admin", "user"},
	}
}
