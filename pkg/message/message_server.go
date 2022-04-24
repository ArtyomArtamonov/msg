package message

import (
	"context"
	"fmt"
	"sync"

	"github.com/ArtyomArtamonov/msg/pkg/auth"
	pb "github.com/ArtyomArtamonov/msg/pkg/message/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageServer struct {
	pb.UnimplementedMessageServiceServer

	mutex            sync.RWMutex
	connectedClients map[string]pb.MessageService_GetMessagesServer
	jwtManager       *auth.JWTManager
}

func NewMessageService(jwtManager *auth.JWTManager) *MessageServer {
	return &MessageServer{
		connectedClients:                  map[string]pb.MessageService_GetMessagesServer{},
		jwtManager:                        jwtManager,
	}
}

func (s *MessageServer) GetMessages(req *emptypb.Empty, srv pb.MessageService_GetMessagesServer) error {
	ctx := srv.Context()
	claims, err := auth.GetAndVerifyClaimsFromContext(ctx, s.jwtManager)
	if err != nil {
		return err
	}

	id := claims.Username

	s.mutex.Lock()
	s.connectedClients[id] = srv
	s.mutex.Unlock()

	logrus.Info("Streaming started with id=", id)
	<-ctx.Done()
	logrus.Info("Streaming ended with id=", id)
	return nil
}

func (s *MessageServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.Status, error) {
	claims, err := auth.GetAndVerifyClaimsFromContext(ctx, s.jwtManager)
	if err != nil {
		return nil, err
	}

	err = s.sendMessage(req.Message, req.To, claims.Username)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user with specified id was not found: %w", err)
	}

	return &pb.Status{
		Success: true,
		Message: fmt.Sprintf("Message was sent to %s", req.To),
	}, nil
}

func (s *MessageServer) sendMessage(message string, to string, from string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client, ok := s.connectedClients[to]
	if !ok {
		return fmt.Errorf("could not read from map with id %s", to)
	}

	msg := pb.MessageResponse{
		From:    from,
		Message: message,
	}
	err := client.Send(&msg)
	logrus.Info("Message sent to id=", to)
	return err
}
