package mocks

import (
	"context"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func (m *UserStoreMock) Save(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return utils.Unwrap[error](args.Get(0))
}

func (m *UserStoreMock) Find(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	return utils.Unwrap[*model.User](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (m *UserStoreMock) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return utils.Unwrap[*model.User](args.Get(0)), utils.Unwrap[error](args.Get(1))
}
