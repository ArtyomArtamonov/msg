package api

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	pb "github.com/ArtyomArtamonov/msg/pkg/api/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageService struct {
	pb.UnimplementedMessageServiceServer

	clients sync.Map
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) GetMessages(req *emptypb.Empty, srv pb.MessageService_GetMessagesServer) error {
	id := s.getId()
	s.clients.Store(id, srv)

	logrus.Info("Streaming started with id=", id)

	<-srv.Context().Done()

	logrus.Info("Streaming ended with id=", id)

	return nil
}

func (s *MessageService) SendMessage(ctx context.Context, req *pb.Message) (*pb.Status, error) {
	id, err := strconv.Atoi(req.To)
	if err != nil {
		return nil, err
	}

	err = s.sendMessage(req.Message, id)
	if err != nil {
		return nil, err
	}

	return &pb.Status{
		Success: true,
		Message: fmt.Sprintf("Message was sent to %s", req.To),
	}, nil
}

func (s MessageService) sendMessage(message string, id int) error {
	client, ok := s.clients.Load(id)
	if !ok {
		return fmt.Errorf("could not read from map with id %d", id)
	}

	msg := pb.Message{
		From:    "Unknown yet (developing)",
		To:      strconv.Itoa(id),
		Message: message,
	}
	err := client.(pb.MessageService_GetMessagesServer).Send(&msg)
	logrus.Info("Message sent to id=", id)
	return err
}

func (s MessageService) getId() int {
	id := 0
	s.clients.Range(func(_, _ any) bool {
		id++
		return true
	})
	return id
}
