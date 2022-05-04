package mocks

import (
	"context"

	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

type ServerStreamMock struct {
	mock.Mock
}

func (m *ServerStreamMock) SetHeader(md metadata.MD) error {
	args := m.Called(md)
	return utils.Unwrap[error](args.Get(0))
}

func (m *ServerStreamMock) SendHeader(md metadata.MD) error {
	args := m.Called(md)
	return utils.Unwrap[error](args.Get(0))
}

func (m *ServerStreamMock) SetTrailer(md metadata.MD) {
	m.Called(md)
}

func (m *ServerStreamMock) Context() context.Context {
	args := m.Called()
	return utils.Unwrap[context.Context](args.Get(0))
}

func (m *ServerStreamMock) SendMsg(message interface{}) error {
	args := m.Called(message)
	return utils.Unwrap[error](args.Get(0))
}

func (m *ServerStreamMock) RecvMsg(message interface{}) error {
	args := m.Called(message)
	return utils.Unwrap[error](args.Get(0))
}
