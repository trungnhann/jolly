package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	"jolly/backend/common/event"
	"jolly/backend/common/log"
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

	if corrID := log.CorrelationIDFromContext(ctx); corrID != "" {
		middleware.SetCorrelationID(corrID, msg)
	}

	return p.pub.Publish(topic, msg)
}
