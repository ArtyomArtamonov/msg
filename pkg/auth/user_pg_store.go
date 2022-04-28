package auth

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(dbname, username, password string) *PostgresUserStore {
	connectionString := fmt.Sprintf(
		"host=database port=5432 sslmode=disable dbname=%s user=%s password=%s",
		dbname, username, password)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		logrus.Fatal(err)
	}

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

func (s *PostgresUserStore) Find(username string) (*User, error) {
	rows := s.db.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username=$1", username)

	user := &User{}
	if err := rows.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}

	return user, nil
}
