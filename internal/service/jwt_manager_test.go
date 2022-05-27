package service

import (
	"testing"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/mocks"
	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const tokenDuration = time.Minute * 10

func TestJWTManager_GenerateSuccess(t *testing.T) {
	setupTest()

	// Generate is also tested in TestJWTManager_GetAndVerifyClaimsSuccess

	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        tokenDuration,
		refreshTokenDuration: tokenDuration,
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
	setupTest()

	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        tokenDuration,
		refreshTokenDuration: tokenDuration,
	}

	token := "some_invalid_access_token"

	res, err := jwtManager.Verify(token)

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "invalid token")
}

func TestJWTManager_VerifyFailsIfTokenExpired(t *testing.T) {
	setupTest()

	// Set old time so generated token is alredy expired now
	utils.MockNow(time.Date(2020, 01, 01, 00, 00, 00, 00, time.Local))

	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        tokenDuration,
		refreshTokenDuration: tokenDuration,
	}
	user, _ := model.NewUser(
		"admin",
		"admin",
		model.ADMIN_ROLE,
	)
	tokenPair, err := jwtManager.Generate(user)

	res, err := jwtManager.Verify(tokenPair.JwtToken)

	assert.Nil(t, res)
	assert.ErrorContains(t, err, "is expired by")
}

func TestJWTManager_VerifySuccess(t *testing.T) {
	setupTest()

	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        tokenDuration,
		refreshTokenDuration: tokenDuration,
	}
	user, _ := model.NewUser(
		"admin",
		"admin",
		model.ADMIN_ROLE,
	)
	tokenPair, err := jwtManager.Generate(user)

	res, err := jwtManager.Verify(tokenPair.JwtToken)

	assert.Equal(
		t,
		&model.UserClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: utils.Now().Add(jwtManager.tokenDuration).Unix(),
				IssuedAt:  utils.Now().Unix(),
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
	setupTest()

	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        tokenDuration,
		refreshTokenDuration: tokenDuration,
	}

	userId := uuid.New()

	res := jwtManager.NewRefreshToken(userId)

	assert.NotEmpty(t, res.Token)
	assert.Equal(t, userId, res.UserId)
	assert.Equal(t, utils.Now().Add(jwtManager.refreshTokenDuration), res.ExpiresAt)
	assert.Equal(t, utils.Now(), res.IssuedAt)
}

func TestJWTManager_GetAndVerifyClaimsFailsIfNoMetadata(t *testing.T) {
	setupTest()

	jwtManager := &JWTManager{}
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(nil)

	res, err := jwtManager.GetAndVerifyClaims(contextMock)

	assert.Nil(t, res)
	assert.ErrorIs(t, err, status.Error(codes.Unauthenticated, "metadata is not provided"))
}

func TestJWTManager_GetAndVerifyClaimsFailsIfNoToken(t *testing.T) {
	setupTest()

	jwtManager := &JWTManager{}
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		// empty metadata
	})

	res, err := jwtManager.GetAndVerifyClaims(contextMock)

	assert.Nil(t, res)
	assert.ErrorIs(t, err, status.Error(codes.Unauthenticated, "authorization token is not provided"))
}

func TestJWTManager_GetAndVerifyClaimsFailsIfTokenIsInvalid(t *testing.T) {
	setupTest()

	jwtManager := &JWTManager{}
	token := "some_invalid_access_token"
	contextMock := new(mocks.ContextMock)
	contextMock.On("Value", mock.Anything).Return(metadata.MD{
		"authorization": {token},
	})

	res, err := jwtManager.GetAndVerifyClaims(contextMock)

	assert.Nil(t, res)
	assert.ErrorIs(t, err, status.Errorf(codes.Unauthenticated, "authorization token is invalid: %v", "invalid token: token contains an invalid number of segments"))
}

func TestJWTManager_GetAndVerifyClaimsSuccess(t *testing.T) {
	setupTest()

	// Generate token and verify we can get claims from it

	jwtManager := &JWTManager{
		secretKey:            "some_key",
		tokenDuration:        tokenDuration,
		refreshTokenDuration: tokenDuration,
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
				ExpiresAt: utils.Now().Add(jwtManager.tokenDuration).Unix(),
				IssuedAt:  utils.Now().Unix(),
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
