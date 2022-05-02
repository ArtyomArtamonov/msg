package mocks

import (
	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type JWTManagerMock struct {
	mock.Mock
}

func (mock *JWTManagerMock) Generate(user *model.User) (*model.TokenPair, error) {
	args := mock.Called(user)
	return utils.Unwrap[*model.TokenPair](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (mock *JWTManagerMock) Verify(accessToken string) (*model.UserClaims, error) {
	args := mock.Called(accessToken)
	return utils.Unwrap[*model.UserClaims](args.Get(0)), utils.Unwrap[error](args.Get(1))
}

func (mock *JWTManagerMock) NewRefreshToken(userId uuid.UUID) *model.RefreshToken {
	args := mock.Called(userId)
	return utils.Unwrap[*model.RefreshToken](args.Get(0))
}
