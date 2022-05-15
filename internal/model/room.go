package model

import (
	"time"

	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
)

type Room struct {
	Id              uuid.UUID
	Name            string
	CreatedAt       time.Time
	Users           []string
	DialogRoom      bool
	LastMessageTime time.Time
}

func NewRoom(name string, dialogRoom bool, users ...string) *Room {
	return &Room{
		Id:              uuid.New(),
		Name:            name,
		CreatedAt:       utils.Now(),
		Users:           users,
		DialogRoom:      dialogRoom,
		LastMessageTime: utils.Now(),
	}
}

func (r *Room) PbRoom() *pb.Room {
	return &pb.Room{
		Id:              r.Id.String(),
		Name:            r.Name,
		CreatedAt:       timestamppb.New(r.CreatedAt),
		DialogRoom:      r.DialogRoom,
		LastMessageTime: timestamppb.New(r.LastMessageTime),
	}
}
