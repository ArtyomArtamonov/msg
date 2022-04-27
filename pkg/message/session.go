package message

import (
	"time"

	pb "github.com/ArtyomArtamonov/msg/pkg/message/proto"
)

type Session struct {
	id         string
	connection pb.MessageService_GetMessagesServer
	expires    time.Duration
	done       chan<- struct{}
}
