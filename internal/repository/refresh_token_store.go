package repository

import (
	"context"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenStore interface {
	Add(ctx context.Context, token *model.RefreshToken) error
	Delete(ctx context.Context, token uuid.UUID) error
	Get(ctx context.Context, token uuid.UUID) (*model.RefreshToken, error)
}

type RefreshTokenPostgresStore struct {
	db *pgxpool.Pool
}

func NewRefreshTokenPostgresStore(db *pgxpool.Pool) *RefreshTokenPostgresStore {
	return &RefreshTokenPostgresStore{
		db: db,
	}
}

func (s *RefreshTokenPostgresStore) Add(ctx context.Context, token *model.RefreshToken) error {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, "INSERT INTO refresh_tokens(token, user_id, expires_at, issued_at) VALUES($1, $2, $3, $4)",
		token.Token, token.UserId, token.ExpiresAt, token.IssuedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *RefreshTokenPostgresStore) Delete(ctx context.Context, token uuid.UUID) error {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, "DELETE FROM refresh_tokens WHERE token=$1", token)

	if err != nil {
		return err
	}

	return nil
}

func (s *RefreshTokenPostgresStore) Get(ctx context.Context, token uuid.UUID) (*model.RefreshToken, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	row := conn.QueryRow(ctx, "SELECT * FROM refresh_tokens WHERE token=$1", token)

	refreshToken := new(model.RefreshToken)
	if err := row.Scan(&refreshToken.Token, &refreshToken.UserId, &refreshToken.ExpiresAt, &refreshToken.IssuedAt); err != nil {
		return nil, err
	}

	return refreshToken, nil
}
