package mocks

import (
	"time"

	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/stretchr/testify/mock"
)

type ContextMock struct {
	mock.Mock
}

func (m *ContextMock) Done() <-chan struct{} {
	args := m.Called()
	return utils.Unwrap[<-chan struct{}](args.Get(0))
}

func (m *ContextMock) Err() error {
	args := m.Called()
	return utils.Unwrap[error](args.Get(0))
}

func (m *ContextMock) Deadline() (deadline time.Time, ok bool) {
	args := m.Called()
	return utils.Unwrap[time.Time](args.Get(0)), utils.Unwrap[bool](args.Get(1))
}

func (m *ContextMock) Value(key any) any {
	args := m.Called(key)
	return utils.Unwrap[any](args.Get(0))
}
