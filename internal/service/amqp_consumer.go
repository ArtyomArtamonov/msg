package service

import (
	"github.com/ArtyomArtamonov/msg/internal/repository"
	pb "github.com/ArtyomArtamonov/msg/internal/server/msg-proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type AMQPConsumer interface {
	Consume()
}

type RabbitMQConsumer struct {
	Channel      *amqp.Channel
	SessionStore repository.SessionStore
	Queue        *amqp.Queue
}

func NewRabbitMQConsumer(channel *amqp.Channel, sessionStore repository.SessionStore) *RabbitMQConsumer {
	queue, err := channel.QueueDeclare(
		uuid.New().String(), // channelname
		false,               // durable
		false,               // delete when unused
		true,                // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		logrus.Fatalf("could not create queue: %s", err.Error())
	}

	err = channel.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		"amq.fanout", // exchange
		false,
		nil,
	)
	if err != nil {
		logrus.Fatalf("could not bind exchange to queue: %s", err.Error())
	}

	return &RabbitMQConsumer{
		Channel:      channel,
		SessionStore: sessionStore,
		Queue:        &queue,
	}
}

func (c *RabbitMQConsumer) Consume() {
	ch, err := c.Channel.Consume(
		c.Queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		logrus.Fatalf("could not create channel: %s", err.Error())
	}

	for {
		select {
		case delivery := <-ch:
			var messageDelivery pb.MessageDelivery
			err := proto.Unmarshal(delivery.Body, &messageDelivery)
			if err != nil {
				logrus.Errorf("could not unmarshal delivery: %v", err)
				continue
			}

			response := &pb.MessageStreamResponse{
				Message: messageDelivery.Message,
			}

			for _, id := range messageDelivery.UserIds {
				// Do not send message to sender
				if messageDelivery.Message.UserId == id {
					continue
				}

				id, err := uuid.Parse(id)
				if err != nil {
					logrus.Warningf("could not parse id: %v", err)
					continue
				}
				err = c.SessionStore.Send(id, response)
				if err != nil {
					continue
				}
			}
		}
	}
}
