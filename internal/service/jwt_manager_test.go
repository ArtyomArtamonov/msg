package service

import (
	"testing"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/mocks"
	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestJWTManager_GenerateSuccess(t *testing.T) {
	// Generate is also tested in TestJWTManager_GetAndVerifyClaimsSuccess

	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        10,
		refreshTokenDuration: 10,
	}

	user, _ := model.NewUser(
		"admin",
		"admin",
		model.ADMIN_ROLE,
	)

	res, err := jwtManager.Generate(user)

	assert.NotEmpty(t, res.JwtToken)
	assert.NotNil(t, res.RefreshToken)
	assert.Nil(t, err)
}

func TestJWTManager_VerifyFailsIfInvalidToken(t *testing.T) {
	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        10,
		refreshTokenDuration: 10,
	}

	token := "some_invalid_access_token"

	res, err := jwtManager.Verify(token)

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "invalid token")
}

func TestJWTManager_VerifySuccess(t *testing.T) {
	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        10,
		refreshTokenDuration: 10,
	}
	user, _ := model.NewUser(
		"admin",
		"admin",
		model.ADMIN_ROLE,
	)
	tokenPair, err := jwtManager.Generate(user)

	res, err := jwtManager.Verify(tokenPair.JwtToken)

	assert.NotEmpty(
		t,
		&model.UserClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(jwtManager.tokenDuration).Unix(),
				IssuedAt:  time.Now().Unix(),
				Id:        user.Id.String(),
				Subject:   user.Id.String(),
			},
			Username: user.Username,
			Role:     user.Role,
		},
		res,
	)
	assert.Nil(t, err)
}

func TestJWTManager_NewRefreshTokenSuccess(t *testing.T) {
	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        10,
		refreshTokenDuration: 10,
	}

	userId := uuid.New()

	res := jwtManager.NewRefreshToken(userId)

	assert.NotEmpty(t, res.Token)
	assert.Equal(t, userId, res.UserId)
	assert.Equal(t, time.Now().Add(jwtManager.refreshTokenDuration), res.ExpiresAt)
	assert.Equal(t, time.Now(), res.IssuedAt)
}

func TestJWTManager_GetAndVerifyClaimsFailsIfNoMetadata(t *testing.T) {
	jwtManager := &JWTManager{}
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(nil)

	res, err := jwtManager.GetAndVerifyClaims(contextMock)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Unauthenticated, "metadata is not provided"), err)
}

func TestJWTManager_GetAndVerifyClaimsFailsIfNoToken(t *testing.T) {
	jwtManager := &JWTManager{}
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		// empty metadata
	})

	res, err := jwtManager.GetAndVerifyClaims(contextMock)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Error(codes.Unauthenticated, "authorization token is not provided"), err)
}

func TestJWTManager_GetAndVerifyClaimsFailsIfTokenIsInvalid(t *testing.T) {
	jwtManager := &JWTManager{}
	token := "some_invalid_access_token"
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		"authorization": {token},
	})

	res, err := jwtManager.GetAndVerifyClaims(contextMock)

	assert.Nil(t, res)
	assert.ErrorIs(t, status.Errorf(codes.Unauthenticated, "authorization token is invalid: %v", "invalid token: token contains an invalid number of segments"), err)
}

func TestJWTManager_GetAndVerifyClaimsSuccess(t *testing.T) {
	// Generate token and verify we can get claims from it

	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        10,
		refreshTokenDuration: 10,
	}

	user, _ := model.NewUser(
		"admin",
		"admin",
		model.ADMIN_ROLE,
	)
	tokenPair, err := jwtManager.Generate(user)

	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		"authorization": {tokenPair.JwtToken},
	})

	res, err := jwtManager.GetAndVerifyClaims(contextMock)

	assert.Equal(
		t,
		&model.UserClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(jwtManager.tokenDuration).Unix(),
				IssuedAt:  time.Now().Unix(),
				Id:        user.Id.String(),
				Subject:   user.Id.String(),
			},
			Username: user.Username,
			Role:     user.Role,
		},
		res,
	)
	assert.Nil(t, err)
}
