package event

import (
	"context"
	"time"
)

type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type Publisher interface {
	Publish(ctx context.Context, topic string, event DomainEvent) error
}
