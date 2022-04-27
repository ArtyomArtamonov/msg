package message

import (
	"context"
	"fmt"
	"time"

	"github.com/ArtyomArtamonov/msg/pkg/auth"
	pb "github.com/ArtyomArtamonov/msg/pkg/message/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageServer struct {
	pb.UnimplementedMessageServiceServer

	sessionStore SessionStore
	jwtManager   *auth.JWTManager
}

func NewMessageService(jwtManager *auth.JWTManager, sessionStore SessionStore) *MessageServer {
	return &MessageServer{
		sessionStore: sessionStore,
		jwtManager:   jwtManager,
	}
}

func (s *MessageServer) GetMessages(req *emptypb.Empty, srv pb.MessageService_GetMessagesServer) error {
	ctx := srv.Context()
	claims, err := auth.GetAndVerifyClaimsFromContext(ctx, s.jwtManager)
	if err != nil {
		return err
	}

	id := claims.Username

	done := make(chan struct{})
	session := Session{
		connection: srv,
		id:         claims.Username,
		expires:    time.Duration(claims.ExpiresAt),
		done:       done,
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
	msg := pb.MessageResponse{
		From:    from,
		Message: message,
	}
	err := s.sessionStore.Send(to, &msg)
	status, ok := status.FromError(err)
	if ok {
		if status.Code() == codes.Unavailable {
			// Client is not connected, no session created. Send push notification
			logrus.Warn("User was not connected. Sending PUSH notification")
		} else if status.Code() == codes.Unauthenticated {
			// JWT token is bad. Remove session
			s.sessionStore.Delete(to)
		}
	} else {
		// Could not send message to client due to unknown error
		logrus.Error(err)
	}

	// TODO: Save message to database
	logrus.Info("Message sent to id=", to)
	return err
}
