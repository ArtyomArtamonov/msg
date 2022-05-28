package repository

import (
	"context"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserStore interface {
	Save(ctx context.Context, user *model.User) error
	Find(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type PostgresUserStore struct {
	db *sqlx.DB
}

func NewPostgresUserStore(db *sqlx.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) Save(ctx context.Context, user *model.User) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO users(id, username, password_hash, role) VALUES($1,$2,$3,$4) RETURNING id",
		user.Id, user.Username, user.PasswordHash, user.Role)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) Find(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := new(model.User)
	err := s.db.GetContext(ctx, user, "SELECT id, username, password_hash, role FROM users WHERE id=$1", id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user := new(model.User)
	err := s.db.GetContext(ctx, user, "SELECT id, username, password_hash, role FROM users WHERE username=$1", username)

	if err != nil {
		return nil, err
	}

	return user, nil
}
