package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTManager struct {
	secretKey            string
	tokenDuration        time.Duration
	refreshTokenDuration time.Duration
}

type TokenPair struct {
	JwtToken     string
	RefreshToken *RefreshToken
}

type RefreshToken struct {
	Token     uuid.UUID
	UserId    uuid.UUID
	ExpiresAt time.Time
	IssuedAt  time.Time
}

type UserClaims struct {
	jwt.StandardClaims

	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewJWTManager(secretKey string, tokenDuration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:            secretKey,
		tokenDuration:        tokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (m *JWTManager) Generate(user *User) (*TokenPair, error) {
	claims := UserClaims{
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

	tokenPair := &TokenPair{
		JwtToken:     tokenSigned,
		RefreshToken: m.NewRefreshToken(user.Id),
	}
	return tokenPair, nil
}

func (m *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected token signing method")
		}

		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (m *JWTManager) NewRefreshToken(userId uuid.UUID) *RefreshToken {
	token := &RefreshToken{
		Token:     uuid.New(),
		UserId:    userId,
		ExpiresAt: time.Now().Add(m.refreshTokenDuration),
		IssuedAt:  time.Now(),
	}
	return token
}
