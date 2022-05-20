package server

import (
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

	done := make(chan error)
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
	
	var doneErr error
	
	select {
	case doneErr = <-done:
	case <-ctx.Done():
	}
	s.sessionStore.Delete(id)
	
	logrus.Info("Streaming ended with id=", id)
	
	return doneErr
}
