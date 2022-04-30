package message

import (
	"time"

	pb "github.com/ArtyomArtamonov/msg/internal/message/proto"
	"github.com/google/uuid"
)

type Session struct {
	id         uuid.UUID
	connection pb.MessageService_GetMessagesServer
	expires    time.Duration
	done       chan<- struct{}
}
