package server

import (
	"context"
	"fmt"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"

	"github.com/ArtyomArtamonov/msg/internal/repository"
	"github.com/ArtyomArtamonov/msg/internal/service"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageServer struct {
	pb.UnimplementedMessageServiceServer

	sessionStore repository.SessionStore
	jwtManager   *service.JWTManager
}

func NewMessageServer(jwtManager *service.JWTManager, sessionStore repository.SessionStore) *MessageServer {
	return &MessageServer{
		sessionStore: sessionStore,
		jwtManager:   jwtManager,
	}
}

func (s *MessageServer) GetMessages(req *emptypb.Empty, srv pb.MessageService_GetMessagesServer) error {
	ctx := srv.Context()
	claims, err := s.jwtManager.GetAndVerifyClaims(ctx)
	if err != nil {
		return err
	}

	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return status.Error(codes.InvalidArgument, "could not parse uuid")
	}

	done := make(chan struct{})
	session := model.Session{
		Connection: srv,
		Id:         id,
		Expires:    time.Duration(claims.ExpiresAt),
		Done:       done,
	}
	err = s.sessionStore.Add(&session)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	logrus.Info("Streaming started with id=", id)
	select {
	case <-done:
		return status.Error(codes.Unauthenticated, "JWT is exipred")
	case <-ctx.Done():
	}
	err = s.sessionStore.Delete(id)
	if err != nil {
		logrus.Warn(err)
	}
	logrus.Info("Streaming ended with id=", id)
	return nil
}

func (s *MessageServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageRequestStatus, error) {
	claims, err := s.jwtManager.GetAndVerifyClaims(ctx)
	if err != nil {
		return nil, err
	}

	err = s.sendMessage(req.Message, req.To, claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user with specified id was not found: %v", err)
	}

	return &pb.MessageRequestStatus{
		Success: true,
		Message: fmt.Sprintf("Message was sent to %s", req.To),
	}, nil
}

func (s *MessageServer) sendMessage(message string, to string, from string) error {
	if _, err := uuid.Parse(from); err != nil {
		return fmt.Errorf("could not parse uuid from %s", from)
	}

	msg := pb.MessageResponse{
		From:    from,
		Message: message,
	}

	idTo, err := uuid.Parse(to)
	if err != nil {
		return fmt.Errorf("could not parse uuid to %s", to)
	}

	err = s.sessionStore.Send(idTo, &msg)
	status, ok := status.FromError(err)
	if ok {
		if status.Code() == codes.Unavailable {
			// Client is not connected, no session created. Send push notification
			logrus.Warn("User was not connected. Sending PUSH notification")
		} else if status.Code() == codes.Unauthenticated {
			// JWT token is bad. Remove session
			s.sessionStore.Delete(idTo)
		}
	} else {
		// Could not send message to client due to unknown error
		logrus.Error(err)
	}

	// TODO: Save message to databasepa
	logrus.Info("Message sent to id=", to)
	return err
}
