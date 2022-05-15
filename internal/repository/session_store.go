package repository

import (
	"sync"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SessionStore interface {
	Add(*model.Session) error
	Send(uuid.UUID, *pb.MessageResponse) error
	Delete(id uuid.UUID) error
}

type InMemorySessionStore struct {
	mutex    sync.Mutex
	sessions map[uuid.UUID]*model.Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		mutex:    sync.Mutex{},
		sessions: make(map[uuid.UUID]*model.Session),
	}
}

func (s *InMemorySessionStore) Add(session *model.Session) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.sessions[session.Id] = session

	return nil
}

func (s *InMemorySessionStore) Send(id uuid.UUID, message *pb.MessageResponse) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return status.Errorf(codes.Unavailable, "User %s is not connected to session", id)
	}

	if time.Duration(utils.Now().Unix()) >= session.Expires {
		return status.Errorf(codes.Unauthenticated, "JWT is expired")
	}

	err := session.Connection.Send(message)

	return err
}

func (s *InMemorySessionStore) Delete(id uuid.UUID) error {
	s.mutex.Lock()

	session, ok := s.sessions[id]
	if !ok {
		return status.Error(codes.NotFound, "User with specified id is not present")
	}

	s.mutex.Unlock()
	delete(s.sessions, id)

	session.Done <- struct{}{}
	return nil
}
