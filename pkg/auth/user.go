package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ADMIN_ROLE = "admin"
	USER_ROLE  = "user"
)

type User struct {
	Id           uuid.UUID
	Username     string
	PasswordHash string
	Role         string
}

func NewUser(username, password, role string) (*User, error) {
	hashedPassword, err := hash(password)
	if err != nil {
		return nil, err
	}

	id := uuid.New()

	return &User{
		Id:           id,
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         role,
	}, nil
}

func (u *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) Clone() *User {
	return &User{
		Id:           u.Id,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
	}
}

func hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
