package mocks

import (
	"context"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MessageStoreMock struct {
	mock.Mock
}

func (m *MessageStoreMock) ListMessages(ctx context.Context, id uuid.UUID, createdAt time.Time, pageSize int) ([]model.Message, error) {
	args := m.Called(ctx, id, pageSize)
	return utils.Unwrap[[]model.Message](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *MessageStoreMock) ListMessagesFirst(ctx context.Context, id uuid.UUID, pageSize int) ([]model.Message, error) {
	args := m.Called(ctx, id, pageSize)
	return utils.Unwrap[[]model.Message](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *MessageStoreMock) SendMessage(ctx context.Context, message *model.Message) (error) {
	args := m.Called(ctx, message)
	return utils.Unwrap[error](args.Get(0))
}
