package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthServer_RegisterFailsIfUserAlreadyExists(t *testing.T) {
	setupTest()

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

func TestAuthServer_RegisterFailsIfUsernameTooLong(t *testing.T) {
	setupTest()

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

func TestAuthServer_RegisterFailsIfPasswordTooShort(t *testing.T) {
	setupTest()

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

func TestAuthServer_RegisterFailsIfJWTFails(t *testing.T) {
	setupTest()

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
	assert.ErrorIs(t, status.Error(codes.Internal, "could not generate token pair"), err)
}

func TestAuthServer_RegisterFailsIfUserStoreFails(t *testing.T) {
	setupTest()

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

func TestAuthServer_RegisterFailsIfRefreshTokenStoreFails(t *testing.T) {
	setupTest()

	expectedUsername := "some_user"
	expectedPassword := "123456"
	expectedError := errors.New("some_error")

	userStoreMock.On("FindByUsername", expectedUsername).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(&model.TokenPair{
		JwtToken: "",
		RefreshToken: &model.RefreshToken{
			Token:     uuid.New(),
			UserId:    uuid.New(),
			ExpiresAt: utils.Now(),
			IssuedAt:  utils.Now(),
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
	assert.ErrorContains(t, err, "could not save refresh token to database")
}

func TestAuthServer_RegisterSuccess(t *testing.T) {
	setupTest()

	expectedUsername := "some_user"
	expectedPassword := "123456"
	expectedTokenPair := &model.TokenPair{
		JwtToken: "",
		RefreshToken: &model.RefreshToken{
			Token:     uuid.New(),
			UserId:    uuid.New(),
			ExpiresAt: utils.Now(),
			IssuedAt:  utils.Now(),
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
			Token: &pb.Token{
				AccessToken:  expectedTokenPair.JwtToken,
				RefreshToken: expectedTokenPair.RefreshToken.Token.String(),
			},
		},
		res,
	)
	assert.Nil(t, err)
}

func TestAuthServer_LoginFailsIfUserDoesNotExist(t *testing.T) {
	setupTest()

	expectedUsername := "some_user"

	userStoreMock.On("FindByUsername", expectedUsername).Return(nil, nil)

	res, err := authServer.Login(
		context.TODO(),
		&pb.LoginRequest{
			Username: expectedUsername,
			Password: "",
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.NotFound, "incorrect username or password"), err)
}

func TestAuthServer_LoginFailsIfPasswordIsNotCorrect(t *testing.T) {
	setupTest()

	expectedUsername := "some_user"
	password := "123456"

	user, err := model.NewUser(
		expectedUsername,
		password,
		model.USER_ROLE,
	)
	userStoreMock.On("FindByUsername", expectedUsername).Return(user, nil)

	res, err := authServer.Login(
		context.TODO(),
		&pb.LoginRequest{
			Username: expectedUsername,
			Password: "incorrect_password",
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.NotFound, "incorrect username or password"), err)
}

func TestAuthServer_LoginFailsIfJWTFails(t *testing.T) {
	setupTest()

	expectedUsername := "some_user"
	password := "123456"
	expectedError := errors.New("some_error")

	user, err := model.NewUser(
		expectedUsername,
		password,
		model.USER_ROLE,
	)
	userStoreMock.On("FindByUsername", expectedUsername).Return(user, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(nil, expectedError)

	res, err := authServer.Login(
		context.TODO(),
		&pb.LoginRequest{
			Username: expectedUsername,
			Password: password,
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Internal, "could not generate token pair"), err)
}

func TestAuthServer_LoginFailsIfRefreshTokenStoreFails(t *testing.T) {
	setupTest()

	expectedUsername := "some_user"
	password := "123456"
	expectedError := errors.New("some_error")

	user, err := model.NewUser(
		expectedUsername,
		password,
		model.USER_ROLE,
	)
	userStoreMock.On("FindByUsername", expectedUsername).Return(user, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(&model.TokenPair{
		JwtToken: "",
		RefreshToken: &model.RefreshToken{
			Token:     uuid.New(),
			UserId:    uuid.New(),
			ExpiresAt: utils.Now(),
			IssuedAt:  utils.Now(),
		},
	}, nil)
	refreshTokenStoreMock.On("Add", mock.Anything).Return(expectedError)

	res, err := authServer.Login(
		context.TODO(),
		&pb.LoginRequest{
			Username: expectedUsername,
			Password: password,
		})

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "could not save refresh token to database")
}

func TestAuthServer_LoginSuccess(t *testing.T) {
	setupTest()

	expectedUsername := "some_user"
	password := "123456"
	expectedTokenPair := &model.TokenPair{
		JwtToken: "",
		RefreshToken: &model.RefreshToken{
			Token:     uuid.New(),
			UserId:    uuid.New(),
			ExpiresAt: utils.Now(),
			IssuedAt:  utils.Now(),
		},
	}

	user, err := model.NewUser(
		expectedUsername,
		password,
		model.USER_ROLE,
	)
	userStoreMock.On("FindByUsername", expectedUsername).Return(user, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(expectedTokenPair, nil)
	refreshTokenStoreMock.On("Add", mock.Anything).Return(nil)

	res, err := authServer.Login(
		context.TODO(),
		&pb.LoginRequest{
			Username: expectedUsername,
			Password: password,
		})

	assert.Equal(
		t,
		&pb.TokenResponse{
			Token: &pb.Token{
				AccessToken:  expectedTokenPair.JwtToken,
				RefreshToken: expectedTokenPair.RefreshToken.Token.String(),
			},
		},
		res,
	)
	assert.Nil(t, err)
}

func TestAuthServer_RefreshFailsIfInvalidRefresh(t *testing.T) {
	setupTest()

	refreshToken := "invalid_token"

	res, err := authServer.Refresh(
		context.TODO(),
		&pb.RefreshRequest{
			RefreshToken: refreshToken,
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.InvalidArgument, "could not parse refresh token"), err)
}

func TestAuthServer_RefreshFailsIfRefreshTokenStoreFails(t *testing.T) {
	setupTest()

	refreshTokenUuid := uuid.New()
	expectedError := errors.New("some_error")

	refreshTokenStoreMock.On("Get", refreshTokenUuid).Return(nil, expectedError)

	res, err := authServer.Refresh(
		context.TODO(),
		&pb.RefreshRequest{
			RefreshToken: refreshTokenUuid.String(),
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Unauthenticated, "refresh token does not exists"), err)
}

func TestAuthServer_RefreshFailsIfRefreshTokenExpires(t *testing.T) {
	setupTest()

	refreshTokenUuid := uuid.New()
	refreshToken := &model.RefreshToken{
		Token:     refreshTokenUuid,
		UserId:    uuid.New(),
		ExpiresAt: utils.Now().Add(-time.Second),
		IssuedAt:  utils.Now(),
	}

	refreshTokenStoreMock.On("Get", refreshTokenUuid).Return(refreshToken, nil)
	refreshTokenStoreMock.On("Delete", refreshTokenUuid).Return(nil)

	res, err := authServer.Refresh(
		context.TODO(),
		&pb.RefreshRequest{
			RefreshToken: refreshTokenUuid.String(),
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Unauthenticated, "refresh token is expired"), err)
}

func TestAuthServer_RefreshFailsIfUserIsNotFound(t *testing.T) {
	setupTest()

	refreshTokenUuid := uuid.New()
	userUuid := uuid.New()
	refreshToken := &model.RefreshToken{
		Token:     refreshTokenUuid,
		UserId:    userUuid,
		ExpiresAt: utils.Now(),
		IssuedAt:  utils.Now(),
	}
	expectedError := errors.New("some_error")

	refreshTokenStoreMock.On("Get", refreshTokenUuid).Return(refreshToken, nil)
	userStoreMock.On("Find", userUuid).Return(nil, expectedError)

	res, err := authServer.Refresh(
		context.TODO(),
		&pb.RefreshRequest{
			RefreshToken: refreshTokenUuid.String(),
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Internal, "hmm... this is strange. That could not possibly happen"), err)
}

func TestAuthServer_RefreshFailsIfJWTFails(t *testing.T) {
	setupTest()

	refreshTokenUuid := uuid.New()
	userUuid := uuid.New()
	refreshToken := &model.RefreshToken{
		Token:     refreshTokenUuid,
		UserId:    userUuid,
		ExpiresAt: utils.Now(),
		IssuedAt:  utils.Now(),
	}
	expectedError := errors.New("some_error")

	refreshTokenStoreMock.On("Get", refreshTokenUuid).Return(refreshToken, nil)
	refreshTokenStoreMock.On("Delete", refreshTokenUuid).Return(nil)
	userStoreMock.On("Find", userUuid).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(nil, expectedError)

	res, err := authServer.Refresh(
		context.TODO(),
		&pb.RefreshRequest{
			RefreshToken: refreshTokenUuid.String(),
		})

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Internal, "could not generate token pair"), err)
}

func TestAuthServer_RefreshFailsIfDeletingTokenFails(t *testing.T) {
	setupTest()

	refreshTokenUuid := uuid.New()
	userUuid := uuid.New()
	refreshToken := &model.RefreshToken{
		Token:     refreshTokenUuid,
		UserId:    userUuid,
		ExpiresAt: utils.Now(),
		IssuedAt:  utils.Now(),
	}
	expectedError := errors.New("some_error")

	refreshTokenStoreMock.On("Get", refreshTokenUuid).Return(refreshToken, nil)
	userStoreMock.On("Find", userUuid).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(nil, nil)
	refreshTokenStoreMock.On("Delete", refreshTokenUuid).Return(expectedError)

	res, err := authServer.Refresh(
		context.TODO(),
		&pb.RefreshRequest{
			RefreshToken: refreshTokenUuid.String(),
		})

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "could not delete old refresh token")
}

func TestAuthServer_RefreshFailsIfCreatingTokenFails(t *testing.T) {
	setupTest()

	refreshTokenUuid := uuid.New()
	userUuid := uuid.New()
	tokenPair := &model.TokenPair{
		JwtToken: "",
		RefreshToken: &model.RefreshToken{
			Token:     refreshTokenUuid,
			UserId:    userUuid,
			ExpiresAt: utils.Now(),
			IssuedAt:  utils.Now(),
		},
	}
	expectedError := errors.New("some_error")

	refreshTokenStoreMock.On("Get", refreshTokenUuid).Return(tokenPair.RefreshToken, nil)
	userStoreMock.On("Find", userUuid).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(tokenPair, nil)
	refreshTokenStoreMock.On("Delete", refreshTokenUuid).Return(nil)
	refreshTokenStoreMock.On("Add", mock.Anything).Return(expectedError)

	res, err := authServer.Refresh(
		context.TODO(),
		&pb.RefreshRequest{
			RefreshToken: refreshTokenUuid.String(),
		})

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "could not save new refresh token")
}

func TestAuthServer_RefreshSuccess(t *testing.T) {
	setupTest()

	refreshTokenUuid := uuid.New()
	userUuid := uuid.New()
	tokenPair := &model.TokenPair{
		JwtToken: "",
		RefreshToken: &model.RefreshToken{
			Token:     refreshTokenUuid,
			UserId:    userUuid,
			ExpiresAt: utils.Now(),
			IssuedAt:  utils.Now(),
		},
	}

	refreshTokenStoreMock.On("Get", refreshTokenUuid).Return(tokenPair.RefreshToken, nil)
	userStoreMock.On("Find", userUuid).Return(nil, nil)
	jwtManagerMock.On("Generate", mock.Anything).Return(tokenPair, nil)
	refreshTokenStoreMock.On("Delete", refreshTokenUuid).Return(nil)
	refreshTokenStoreMock.On("Add", mock.Anything).Return(nil)

	res, err := authServer.Refresh(
		context.TODO(),
		&pb.RefreshRequest{
			RefreshToken: tokenPair.RefreshToken.Token.String(),
		})

	assert.Equal(
		t,
		&pb.TokenResponse{
			Token: &pb.Token{
				AccessToken:  tokenPair.JwtToken,
				RefreshToken: tokenPair.RefreshToken.Token.String(),
			},
		},
		res,
	)
	assert.Nil(t, err)
}
