package mocks

import (
	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (m *UserStoreMock) Save(user *model.User) error {
	args := m.Called(user)
	return utils.Unwrap[error](args.Get(0))
}

func (m *UserStoreMock) Find(id uuid.UUID) (*model.User, error) {
	args := m.Called(id)
	return utils.Unwrap[*model.User](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *UserStoreMock) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return utils.Unwrap[*model.User](args.Get(0)), utils.Unwrap[error](args.Get(1))
}
