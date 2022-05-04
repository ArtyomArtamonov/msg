package mocks

import (
	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type RefreshTokenStoreMock struct {
	mock.Mock
}

func (m *RefreshTokenStoreMock) Add(token *model.RefreshToken) error {
	args := m.Called(token)
	return utils.Unwrap[error](args.Get(0))
}

func (m *RefreshTokenStoreMock) Delete(token uuid.UUID) error {
	args := m.Called(token)
	return utils.Unwrap[error](args.Get(0))
}

func (m *RefreshTokenStoreMock) Get(token uuid.UUID) (*model.RefreshToken, error) {
	args := m.Called(token)
	return utils.Unwrap[*model.RefreshToken](args.Get(0)), utils.Unwrap[error](args.Get(1))
}
