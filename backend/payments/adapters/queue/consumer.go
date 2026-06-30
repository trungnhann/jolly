package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"jolly/backend/common/event"
	"jolly/backend/payments/app"
	"jolly/backend/payments/domain"
)

type IncomingOrderCreatedEvent struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Currency   string `json:"currency"`
	TotalCents int64  `json:"total_cents"`
}

type Consumer struct {
	service   *app.Service
	publisher event.Publisher
}

func NewConsumer(service *app.Service, publisher event.Publisher) *Consumer {
	return &Consumer{service: service, publisher: publisher}
}

func (c *Consumer) HandleOrderCreated(msg *message.Message) error {
	var ev IncomingOrderCreatedEvent
	if err := json.Unmarshal(msg.Payload, &ev); err != nil {
		return err
	}

	// req := client.AuthorizePaymentRequest{
	// 	OrderID:     ev.OrderID,
	// 	CustomerID:  ev.CustomerID,
	// 	AmountCents: ev.TotalCents,
	// 	Currency:    ev.Currency,
	// }

	// resp, err := c.service.Authorize(context.Background(), req)
	// if err != nil {
	// 	failEvent := domain.PaymentFailedEvent{
	// 		OrderID:   ev.OrderID,
	// 		Reason:    err.Error(),
	// 		CreatedAt: time.Now().UTC(),
	// 	}
	// 	return c.publisher.Publish(context.Background(), failEvent.EventName(), failEvent)
	// }

	successEvent := domain.PaymentAuthorizedEvent{
		OrderID: ev.OrderID,
		// PaymentID: resp.PaymentID,
		Amount:    ev.TotalCents,
		Currency:  ev.Currency,
		CreatedAt: time.Now().UTC(),
	}
	return c.publisher.Publish(context.Background(), successEvent.EventName(), successEvent)
}
