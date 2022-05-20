package model

import (
	"time"

	pb "github.com/ArtyomArtamonov/msg/internal/server/proto"
	"github.com/google/uuid"
)

type Session struct {
	Id         uuid.UUID
	Connection pb.MessageService_GetMessagesServer
	Expires    time.Duration
	Done       chan<- error
}
