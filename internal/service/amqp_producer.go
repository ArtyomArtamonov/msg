package service

import (
	"google.golang.org/protobuf/proto"
	"github.com/streadway/amqp"
	pb "github.com/ArtyomArtamonov/msg/internal/server/msg-proto"
)

type AMQPProducer interface {
	Produce(*pb.MessageDelivery) error
}

type RabbitMQProducer struct {
	Channel *amqp.Channel
}

func NewRabbitMQManager(channel *amqp.Channel) *RabbitMQProducer {
	return &RabbitMQProducer{
		Channel: channel,
	}
}

func (m *RabbitMQProducer) Produce(delivery *pb.MessageDelivery) error {
	data, err := proto.Marshal(delivery)
	if err != nil {
		return err
	}

	return m.Channel.Publish(
		"amq.fanout", // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/protobuf",
			Body:        data,
		})
}
