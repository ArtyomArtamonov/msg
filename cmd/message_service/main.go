package main

import (
	"fmt"
	"net"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/repository"
	"github.com/ArtyomArtamonov/msg/internal/server"
	"github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/ArtyomArtamonov/msg/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.TraceLevel)

	err := godotenv.Load("../../.env")
	failOnError(err, "Error loading .env file")

	env := server.NewEnv()

	host := env.MESSAGE_HOST
	lis, err := net.Listen("tcp", host)
	failOnError(err, "could not create tcp connection")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@message-broker:5672/",
		env.RABBITMQ_DEFAULT_USER,
		env.RABBITMQ_DEFAULT_PASS))
	failOnError(err, "could not connect to message-broker")

	ch, err := conn.Channel()
	failOnError(err, "could not create channel (message-broker)")
	defer ch.Close()

	grpcServer := createAndPrepareGRPCServer(ch, env)
	
	logrus.Info("Starting grpc server on ", host)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatal(err.Error())
	}
}

func createAndPrepareGRPCServer(ch *amqp.Channel, env *server.Env) *grpc.Server {
	endpoints := server.NewEndpoints()
	endpointRoles := server.NewEndpointRoles(endpoints)

	jwtManager := service.NewJWTManager(
		env.JWT_SECRET,
		time.Minute*time.Duration(env.JWT_DURATION_MIN),
		time.Hour*24*time.Duration(env.REFRESH_DURATION_DAYS),
	)
	sessionStore := repository.NewInMemorySessionStore()
	messageServer := server.NewMessageServer(jwtManager, sessionStore)

	authInterceptor := server.NewAuthInterceptor(jwtManager, endpointRoles)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	proto.RegisterMessageServiceServer(grpcServer, messageServer)
	reflection.Register(grpcServer)

	amqpConsumer := service.NewRabbitMQConsumer(ch, sessionStore)
	go amqpConsumer.Consume()

	return grpcServer
}

func failOnError(err error, text string) {
	if err != nil {
		logrus.Fatalf("%s: %v", text, err)
	}
}
