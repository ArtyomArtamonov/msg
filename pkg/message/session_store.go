package message

import (
	"sync"
	"time"

	pb "github.com/ArtyomArtamonov/msg/pkg/message/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SessionStore interface {
	Add(*Session) error
	Send(string, *pb.MessageResponse) error
	Delete(id string) error
}

type InMemorySessionStore struct {
	mutex    sync.Mutex
	sessions map[string]*Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		mutex:    sync.Mutex{},
		sessions: map[string]*Session{},
	}
}

func (s *InMemorySessionStore) Add(session *Session) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.sessions[session.id] = session

	return nil
}

func (s *InMemorySessionStore) Send(id string, message *pb.MessageResponse) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return status.Errorf(codes.Unavailable, "User %s is not connected to session", id)
	}

	if time.Duration(time.Now().Unix()) >= session.expires {
		return status.Errorf(codes.Unauthenticated, "JWT is expired")
	}

	err := session.connection.Send(message)

	return err
}

func (s *InMemorySessionStore) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return status.Error(codes.NotFound, "User with specified id is not present")
	}

	session.done <- struct{}{}
	delete(s.sessions, id)

	return nil
}
