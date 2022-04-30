package repository

import (
	"database/sql"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
)

type RefreshTokenStore interface {
	Add(token *model.RefreshToken) error
	Delete(token uuid.UUID) error
	Get(token uuid.UUID) (*model.RefreshToken, error)
}

type RefreshTokenPostgresStore struct {
	db *sql.DB
}

func NewRefreshTokenPostgresStore(db *sql.DB) *RefreshTokenPostgresStore {
	return &RefreshTokenPostgresStore{
		db: db,
	}
}

func (s *RefreshTokenPostgresStore) Add(token *model.RefreshToken) error {
	_, err := s.db.Exec("INSERT INTO refresh_tokens(token, user_id, expires_at, issued_at) VALUES($1, $2, $3, $4)",
		token.Token, token.UserId, token.ExpiresAt, token.IssuedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *RefreshTokenPostgresStore) Delete(token uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM refresh_tokens WHERE token=$1", token)

	if err != nil {
		return err
	}

	return nil
}

func (s *RefreshTokenPostgresStore) Get(token uuid.UUID) (*model.RefreshToken, error) {
	row := s.db.QueryRow("SELECT * FROM refresh_tokens WHERE token=$1", token)

	if err := row.Err(); err != nil {
		return nil, err
	}

	refreshToken := new(model.RefreshToken)
	if err := row.Scan(&refreshToken.Token, &refreshToken.UserId, &refreshToken.ExpiresAt, &refreshToken.IssuedAt); err != nil {
		return nil, err
	}

	return refreshToken, nil
}
