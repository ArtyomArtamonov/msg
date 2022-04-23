package main

import (
	"log"
	"net"

	api "github.com/ArtyomArtamonov/msg/pkg/api"
	pb "github.com/ArtyomArtamonov/msg/pkg/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err.Error())
	}
	grpcServer := grpc.NewServer()

	server := api.NewMessageService()
	pb.RegisterMessageServiceServer(grpcServer, server)
	
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
