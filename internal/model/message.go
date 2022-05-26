package model

import (
	"time"

	pb "github.com/ArtyomArtamonov/msg/internal/server/msg-proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	Id        uuid.UUID `db:"id"`
	UserId    uuid.UUID `db:"user_id"`
	RoomId    uuid.UUID `db:"room_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}

func NewMessage(userId, roomId uuid.UUID, text string) *Message {
	return &Message{
		Id:        uuid.New(),
		UserId:    userId,
		RoomId:    roomId,
		Text:      text,
		CreatedAt: time.Now(),
	}
}

func (m *Message) ToPbMessage() *pb.Message {
	return &pb.Message{
		Id:        m.Id.String(),
		RoomId:    m.RoomId.String(),
		UserId:    m.UserId.String(),
		Text:      m.Text,
		CreatedAt: timestamppb.New(m.CreatedAt),
	}
}
