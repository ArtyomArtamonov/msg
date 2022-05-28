package mocks

import (
	"context"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type RoomStoreMock struct {
	mock.Mock
}

func (m *RoomStoreMock) Add(ctx context.Context, room *model.Room) error {
	args := m.Called(ctx, room)
	return utils.Unwrap[error](args.Get(0))
}

func (m *RoomStoreMock) AddAndSendMessage(ctx context.Context, room *model.Room, message *model.Message) (*model.Room, error) {
	args := m.Called(ctx, room, message)
	return utils.Unwrap[*model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) Get(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	args := m.Called(ctx, id)
	return utils.Unwrap[*model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) FindDialogRoom(ctx context.Context, userId1, userId2 uuid.UUID) (*model.Room, error) {
	args := m.Called(ctx, userId1, userId2)
	return utils.Unwrap[*model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) UsersInRoom(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	args := m.Called(ctx, id)
	return utils.Unwrap[[]uuid.UUID](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) FindByIds(ctx context.Context, userIds ...uuid.UUID) ([]model.User, error) {
	args := m.Called(ctx, userIds)
	return utils.Unwrap[[]model.User](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) ListRooms(ctx context.Context, userId uuid.UUID, lastMessageDate time.Time, pageSize int) ([]model.Room, error) {
	args := m.Called(ctx, userId, lastMessageDate, pageSize)
	return utils.Unwrap[[]model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *RoomStoreMock) ListRoomsFirst(ctx context.Context, userId uuid.UUID, pageSize int) ([]model.Room, error) {
	args := m.Called(ctx, userId, pageSize)
	return utils.Unwrap[[]model.Room](args.Get(0)), utils.Unwrap[error](args.Get(1))
}
