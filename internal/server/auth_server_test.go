package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/mocks"
	"github.com/ArtyomArtamonov/msg/internal/model"
	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtManagerMock *mocks.JWTManagerMock
var refreshTokenStoreMock *mocks.RefreshTokenStoreMock
var userStoreMock *mocks.UserStoreMock
var authServer *AuthServer

func setup() {
	jwtManagerMock = new(mocks.JWTManagerMock)
	refreshTokenStoreMock = new(mocks.RefreshTokenStoreMock)
	userStoreMock = new(mocks.UserStoreMock)
	authServer = &AuthServer{
		userStore:         userStoreMock,
		refreshTokenStore: refreshTokenStoreMock,
		jwtManager:        jwtManagerMock,
	}
}

func TestRegisterFailsIfUserAlreadyExists(t *testing.T) {
	setup()

	expectedUsername := "some_user"

	user, err := model.NewUser(
		expectedUsername,
		"",
		model.USER_ROLE,
	)
	userStoreMock.On("FindByUsername", expectedUsername).Return(user, err)

	res, err := authServer.Register(
		context.TODO(),
		&pb.RegisterRequest{
			Username: expectedUsername,
			Password: "",
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.AlreadyExists, "user already exists"), err)
}

func TestRegisterFailsIfUsernameTooLong(t *testing.T) {
	setup()

	expectedUsername := "some_very_long_user_name"

	userStoreMock.On("FindByUsername", expectedUsername).Return(nil, nil)

	res, err := authServer.Register(
		context.TODO(),
		&pb.RegisterRequest{
			Username: expectedUsername,
			Password: "",
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.InvalidArgument, "username could not be more than 15 characters"), err)
}

func TestRegisterFailsIfPasswordTooShort(t *testing.T) {
	setup()

	expectedUsername := "some_user"
	expectedPassword := "12345"

	userStoreMock.On("FindByUsername", expectedUsername).Return(nil, nil)

	res, err := authServer.Register(
		context.TODO(),
		&pb.RegisterRequest{
			Username: expectedUsername,
			Password: expectedPassword,
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.InvalidArgument, "password could not be less than 6 characters"), err)
}

func TestRegisterFailsIfJWTFails(t *testing.T) {
	setup()

	expectedUsername := "some_user"
	expectedPassword := "123456"
	expectedError := errors.New("some_error")

	userStoreMock.On("FindByUsername", expectedUsername).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(nil, expectedError)

	res, err := authServer.Register(
		context.TODO(),
		&pb.RegisterRequest{
			Username: expectedUsername,
			Password: expectedPassword,
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, expectedError, err)
}

func TestRegisterFailsIfUserStoreFails(t *testing.T) {
	setup()

	expectedUsername := "some_user"
	expectedPassword := "123456"
	expectedError := errors.New("some_error")

	userStoreMock.On("FindByUsername", expectedUsername).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(nil, nil)
	userStoreMock.On("Save", mock.Anything).Return(expectedError)

	res, err := authServer.Register(
		context.TODO(),
		&pb.RegisterRequest{
			Username: expectedUsername,
			Password: expectedPassword,
		})

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "could not save user")
}

func TestRegisterFailsIfRefreshTokenStoreFails(t *testing.T) {
	setup()

	expectedUsername := "some_user"
	expectedPassword := "123456"
	expectedError := errors.New("some_error")

	userStoreMock.On("FindByUsername", expectedUsername).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(&model.TokenPair{
		JwtToken: "",
		RefreshToken: &model.RefreshToken{
			Token:     uuid.New(),
			UserId:    uuid.New(),
			ExpiresAt: time.Now(),
			IssuedAt:  time.Now(),
		},
	}, nil)
	userStoreMock.On("Save", mock.Anything).Return(nil)
	refreshTokenStoreMock.On("Add", mock.Anything).Return(expectedError)

	res, err := authServer.Register(
		context.TODO(),
		&pb.RegisterRequest{
			Username: expectedUsername,
			Password: expectedPassword,
		})

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "could not add refresh key to database")
}

func TestRegisterSuccess(t *testing.T) {
	setup()

	expectedUsername := "some_user"
	expectedPassword := "123456"
	expectedTokenPair := &model.TokenPair{
		JwtToken: "",
		RefreshToken: &model.RefreshToken{
			Token:     uuid.New(),
			UserId:    uuid.New(),
			ExpiresAt: time.Now(),
			IssuedAt:  time.Now(),
		},
	}

	userStoreMock.On("FindByUsername", expectedUsername).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(expectedTokenPair, nil)
	userStoreMock.On("Save", mock.Anything).Return(nil)
	refreshTokenStoreMock.On("Add", mock.Anything).Return(nil)

	res, err := authServer.Register(
		context.TODO(),
		&pb.RegisterRequest{
			Username: expectedUsername,
			Password: expectedPassword,
		})

	assert.Equal(
		t,
		&pb.TokenResponse{
			AccessToken:  expectedTokenPair.JwtToken,
			RefreshToken: expectedTokenPair.RefreshToken.Token.String(),
		},
		res,
	)
	assert.Nil(t, err)
}
