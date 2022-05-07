package repository

import (
	"context"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type UserStore interface {
	Save(ctx context.Context, user *model.User) error
	Find(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type PostgresUserStore struct {
	db *pgxpool.Pool
}

func NewPostgresUserStore(db *pgxpool.Pool) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) Save(ctx context.Context, user *model.User) error {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, "INSERT INTO users(id, username, password_hash, role) VALUES($1,$2,$3,$4) RETURNING id",
		user.Id, user.Username, user.PasswordHash, user.Role)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) Find(ctx context.Context, id uuid.UUID) (*model.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	row := conn.QueryRow(ctx, "SELECT id, username, password_hash, role FROM users WHERE id=$1", id)

	user := new(model.User)
	if err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	row := conn.QueryRow(ctx, "SELECT id, username, password_hash, role FROM users WHERE username=$1", username)

	user := new(model.User)
	if err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}

	return user, nil
}
