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
	Token     uuid.UUID `db:"token"`
	UserId    uuid.UUID `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
	IssuedAt  time.Time `db:"issued_at"`
}

type UserClaims struct {
	jwt.StandardClaims

	Username string `json:"username"`
	Role     string `json:"role"`
}
