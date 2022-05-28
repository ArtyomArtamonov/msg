package repository

import (
	"context"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenStore interface {
	Add(ctx context.Context, token *model.RefreshToken) error
	Delete(ctx context.Context, token uuid.UUID) error
	Get(ctx context.Context, token uuid.UUID) (*model.RefreshToken, error)
}

type RefreshTokenPostgresStore struct {
	db *sqlx.DB
}

func NewRefreshTokenPostgresStore(db *sqlx.DB) *RefreshTokenPostgresStore {
	return &RefreshTokenPostgresStore{
		db: db,
	}
}

func (s *RefreshTokenPostgresStore) Add(ctx context.Context, token *model.RefreshToken) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO refresh_tokens(token, user_id, expires_at, issued_at) VALUES($1, $2, $3, $4)",
		token.Token, token.UserId, token.ExpiresAt, token.IssuedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *RefreshTokenPostgresStore) Delete(ctx context.Context, token uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE token=$1", token)

	if err != nil {
		return err
	}

	return nil
}

func (s *RefreshTokenPostgresStore) Get(ctx context.Context, token uuid.UUID) (*model.RefreshToken, error) {
	refreshToken := new(model.RefreshToken)
	err := s.db.GetContext(ctx, refreshToken, "SELECT * FROM refresh_tokens WHERE token=$1", token)

	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}
