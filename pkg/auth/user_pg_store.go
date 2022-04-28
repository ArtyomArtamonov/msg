package auth

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserStore interface {
	Save(user *User) error
	Find(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) Save(user *User) error {
	row := s.db.QueryRow("INSERT INTO users(id, username, password_hash, role) VALUES($1,$2,$3,$4) RETURNING id",
		user.Id, user.Username, user.PasswordHash, user.Role)

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return err
	}

	user.Id = id
	return nil
}

func (s *PostgresUserStore) Find(id uuid.UUID) (*User, error) {
	rows := s.db.QueryRow("SELECT id, username, password_hash, role FROM users WHERE id=$1", id)

	user := &User{}
	if err := rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) FindByUsername(username string) (*User, error) {
	rows := s.db.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username=$1", username)

	user := &User{}
	if err := rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}

	return user, nil
}
