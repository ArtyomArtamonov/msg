package model

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

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
