package api

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	pb "github.com/ArtyomArtamonov/msg/pkg/api/proto"
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
	s.clients.Store(0, srv)
	
	<-srv.Context().Done()

	return nil
}

func (s *MessageService) SendMessage(ctx context.Context, req *pb.Message) (*pb.Status, error) {
	id, err := strconv.Atoi(req.To)
	if err != nil {
		return nil, err
	}

	if id != 0 {
		return nil, fmt.Errorf("id %d is not supported in this version (use id: 0)", id)
	}

	err = s.sendMessage(req.Message, id)
	if err != nil {
		return nil, err
	}

	return &pb.Status{
		Success: true,
	}, nil
}

func (s MessageService) sendMessage(message string, id int) error {
	client, ok := s.clients.Load(id)
	if !ok {
		return fmt.Errorf("could not read from map with id %d", id)
	}

	msg := pb.Response{
		Sender:  "Unknown yet (developing)",
		Message: message,
	}
	err := client.(pb.MessageService_GetMessagesServer).Send(&msg)
	return err
}
