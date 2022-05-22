package repository

import (
	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserStore interface {
	Save(user *model.User) error
	Find(id uuid.UUID) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type PostgresUserStore struct {
	db *sqlx.DB
}

func NewPostgresUserStore(db *sqlx.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) Save(user *model.User) error {
	_, err := s.db.Exec("INSERT INTO users(id, username, password_hash, role) VALUES($1,$2,$3,$4) RETURNING id",
		user.Id, user.Username, user.PasswordHash, user.Role)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) Find(id uuid.UUID) (*model.User, error) {
	user := new(model.User)
	err := s.db.Get(user, "SELECT id, username, password_hash, role FROM users WHERE id=$1", id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) FindByUsername(username string) (*model.User, error) {
	user := new(model.User)
	err := s.db.Get(user, "SELECT id, username, password_hash, role FROM users WHERE username=$1", username)

	if err != nil {
		return nil, err
	}

	return user, nil
}
