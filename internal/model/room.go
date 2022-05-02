package model

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	Id         uuid.UUID
	Name       string
	CreatedAt  time.Time
	Users      []string
	DialogRoom bool
}

func NewRoom(name string, dialogRoom bool, users ...string) *Room {
	return &Room{
		Id:         uuid.New(),
		Name:       name,
		CreatedAt:  time.Now(),
		Users:      users,
		DialogRoom: dialogRoom,
	}
}
