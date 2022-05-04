package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTManagerProtol interface {
	Generate(user *model.User) (*model.TokenPair, error)
	Verify(accessToken string) (*model.UserClaims, error)
	NewRefreshToken(userId uuid.UUID) *model.RefreshToken
	GetAndVerifyClaims(ctx context.Context) (*model.UserClaims, error)
}

type JWTManager struct {
	secretKey            string
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            secretKey,
		tokenDuration:        tokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (m *JWTManager) Generate(user *model.User) (*model.TokenPair, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        user.Id.String(),
			Subject:   user.Id.String(),
		},
		Username: user.Username,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return nil, err
	}

	tokenPair := &model.TokenPair{
		JwtToken:     tokenSigned,
		RefreshToken: m.NewRefreshToken(user.Id),
	}
	return tokenPair, nil
}

func (m *JWTManager) Verify(accessToken string) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&model.UserClaims{},
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(m.secretKey), nil
		})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (m *JWTManager) NewRefreshToken(userId uuid.UUID) *model.RefreshToken {
	token := &model.RefreshToken{
		Token:     uuid.New(),
		UserId:    userId,
		ExpiresAt: time.Now().Add(m.refreshTokenDuration),
		IssuedAt:  time.Now(),
	}
	return token
}

func (m *JWTManager) GetAndVerifyClaims(ctx context.Context) (*model.UserClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := m.Verify(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is invalid: %v", err)
	}
	return claims, nil
}
