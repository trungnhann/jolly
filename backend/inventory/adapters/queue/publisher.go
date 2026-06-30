package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"

	"jolly/backend/common/event"
)

type Publisher struct {
	pub message.Publisher
}

func NewPublisher(pub message.Publisher) *Publisher {
	return &Publisher{pub: pub}
}

func (p *Publisher) Publish(ctx context.Context, topic string, e event.DomainEvent) error {
	payloadBytes, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	msg := message.NewMessage(watermill.NewUUID(), payloadBytes)
	return p.pub.Publish(topic, msg)
}
