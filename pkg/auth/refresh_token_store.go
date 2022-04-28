package auth

import (
	"database/sql"

	"github.com/google/uuid"
)

type RefreshTokenStore interface {
	Add(token *RefreshToken) error
	Delete(userId uuid.UUID) error
	Get(token uuid.UUID) (*RefreshToken, error)
}

type RefreshTokenPostgresStore struct {
	db *sql.DB
}

func NewRefreshTokenPostgresStore(db *sql.DB) *RefreshTokenPostgresStore {
	return &RefreshTokenPostgresStore{
		db: db,
	}
}

func (s *RefreshTokenPostgresStore) Add(token *RefreshToken) error {
	row := s.db.QueryRow("INSERT INTO refresh_tokens(token, user_id, expires_at, issued_at) VALUES($1, $2, $3, $4)",
		token.Token, token.UserId, token.ExpiresAt, token.IssuedAt)

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

func (s *RefreshTokenPostgresStore) Delete(userId uuid.UUID) error {
	row := s.db.QueryRow("DELETE FROM refresh_tokens WHERE token=$1", userId)

	if err := row.Err(); err != nil {
		return err
	}

	return nil
}

func (s *RefreshTokenPostgresStore) Get(token uuid.UUID) (*RefreshToken, error) {
	row := s.db.QueryRow("SELECT * FROM refresh_tokens WHERE token=$1", token)

	if err := row.Err(); err != nil {
		return nil, err
	}

	refreshToken := &RefreshToken{}
	row.Scan(&refreshToken.Token, &refreshToken.UserId, &refreshToken.ExpiresAt, &refreshToken.IssuedAt)

	return refreshToken, nil
}
