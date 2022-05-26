package mocks

import (
	proto "github.com/ArtyomArtamonov/msg/internal/server/msg-proto"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/stretchr/testify/mock"
)

type AMQPProducerMock struct {
	mock.Mock
}

func (m *AMQPProducerMock) Produce(message *proto.MessageDelivery) error {
	args := m.Called(message)
	return utils.Unwrap[error](args.Get(0))
}
