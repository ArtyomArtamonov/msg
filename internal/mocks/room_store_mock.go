package mocks

import (
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type RoomStoreMock struct {
	mock.Mock
}

func (m *RoomStoreMock) Add(room *model.Room) error {
	args := m.Called(room)
	return utils.Unwrap[error](args.Get(0))
}

func (m *RoomStoreMock) Get(id uuid.UUID) (*model.Room, error) {
	args := m.Called(id)
	return utils.Unwrap[*model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) FindDialogRoom(userId1, userId2 uuid.UUID) (*model.Room, error) {
	args := m.Called(userId1, userId2)
	return utils.Unwrap[*model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) UsersInRoom(id uuid.UUID) ([]model.User, error) {
	args := m.Called(id)
	return utils.Unwrap[[]model.User](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) FindByIds(usernames ...string) ([]model.User, error) {
	args := m.Called(usernames)
	return utils.Unwrap[[]model.User](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) ListRooms(userId uuid.UUID, lastMessageDate time.Time, pageSize int) ([]model.Room, error) {
	args := m.Called(userId, lastMessageDate, pageSize)
	return utils.Unwrap[[]model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) ListRoomsFirst(userId uuid.UUID, pageSize int) ([]model.Room, error) {
	args := m.Called(userId, pageSize)
	return utils.Unwrap[[]model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}
