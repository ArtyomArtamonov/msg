package repository

import (
	"database/sql"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserStore interface {
	Save(user *model.User) error
	Find(id uuid.UUID) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
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
	row := s.db.QueryRow("SELECT id, username, password_hash, role FROM users WHERE id=$1", id)

	user := new(model.User)
	if err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) FindByUsername(username string) (*model.User, error) {
	row := s.db.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username=$1", username)

	user := new(model.User)
	if err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}

	return user, nil
}
