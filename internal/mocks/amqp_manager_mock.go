package mocks

import (
	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/stretchr/testify/mock"
)

type AMQPProducerMock struct {
	mock.Mock
}

func (m *AMQPProducerMock) Produce(message *model.Message) error {
	args := m.Called(message)
	return utils.Unwrap[error](args.Get(0))
}
