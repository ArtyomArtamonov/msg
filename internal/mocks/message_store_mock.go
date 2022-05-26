package mocks

import (
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MessageStoreMock struct {
	mock.Mock
}

func (m *MessageStoreMock) ListMessages(id uuid.UUID, createdAt time.Time, pageSize int) ([]model.Message, error) {
	args := m.Called(id, pageSize)
	return utils.Unwrap[[]model.Message](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *MessageStoreMock) ListMessagesFirst(id uuid.UUID, pageSize int) ([]model.Message, error) {
	args := m.Called(id, pageSize)
	return utils.Unwrap[[]model.Message](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *MessageStoreMock) SendMessage(message *model.Message) (error) {
	args := m.Called(message)
	return utils.Unwrap[error](args.Get(0))
}
