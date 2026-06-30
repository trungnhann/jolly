package messaging

import (
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
)

func NewRabbitMQPublisher(logger watermill.LoggerAdapter) (*amqp.Publisher, error) {
	amqpURI := os.Getenv("RABBITMQ_URL")
	if amqpURI == "" {
		amqpURI = "amqp://jolly:jolly@localhost:5672/"
	}

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)

	publisher, err := amqp.NewPublisher(amqpConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("could not create amqp publisher: %w", err)
	}

	return publisher, nil
}

func NewRabbitMQSubscriber(logger watermill.LoggerAdapter) (*amqp.Subscriber, error) {
	amqpURI := os.Getenv("RABBITMQ_URL")
	if amqpURI == "" {
		amqpURI = "amqp://jolly:jolly@localhost:5672/"
	}

	amqpConfig := amqp.NewDurableQueueConfig(amqpURI)

	subscriber, err := amqp.NewSubscriber(amqpConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("could not create amqp subscriber: %w", err)
	}

	return subscriber, nil
}
