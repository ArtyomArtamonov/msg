package model

import (
	"time"

	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
)

type Room struct {
	Id              uuid.UUID   `db:"id"`
	Name            string      `db:"name"`
	CreatedAt       time.Time   `db:"created_at"`
	UserIds         []uuid.UUID `db:"-"`
	DialogRoom      bool        `db:"dialog_room"`
	LastMessageTime time.Time   `db:"last_message_time"`
}

func NewRoom(name string, dialogRoom bool, users ...uuid.UUID) *Room {
	return &Room{
		Id:              uuid.New(),
		Name:            name,
		CreatedAt:       utils.Now(),
		UserIds:         users,
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
