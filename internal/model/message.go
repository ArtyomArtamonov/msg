package model

import (
	"time"

	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	RoomId    uuid.UUID
	Text      string
	CreatedAt time.Time
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
