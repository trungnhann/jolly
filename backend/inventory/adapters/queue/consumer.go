package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"jolly/backend/common/event"
	"jolly/backend/inventory/app/command"
	"jolly/backend/inventory/domain"
)

type IncomingOrderPaidEvent struct {
	OrderID string `json:"order_id"`
	Items   []struct {
		SKU      string `json:"sku"`
		Quantity int    `json:"quantity"`
	} `json:"items"`
}

type Consumer struct {
	commandHandlers *command.Handlers
	publisher       event.Publisher
}

func NewConsumer(commandHandlers *command.Handlers, publisher event.Publisher) *Consumer {
	return &Consumer{commandHandlers: commandHandlers, publisher: publisher}
}

func (c *Consumer) HandleOrderPaid(msg *message.Message) error {
	var ev IncomingOrderPaidEvent
	if err := json.Unmarshal(msg.Payload, &ev); err != nil {
		return err
	}

	cmd := command.Reserve{
		OrderID: ev.OrderID,
		Items: make([]struct {
			SKU      string
			Quantity int
		}, 0, len(ev.Items)),
	}
	for _, item := range ev.Items {
		cmd.Items = append(cmd.Items, struct {
			SKU      string
			Quantity int
		}{
			SKU:      item.SKU,
			Quantity: item.Quantity,
		})
	}

	err := c.commandHandlers.Reserve(context.Background(), cmd)
	if err != nil {
		failEvent := domain.InventoryReservationFailedEvent{
			OrderID:   ev.OrderID,
			Reason:    err.Error(),
			CreatedAt: time.Now().UTC(),
		}
		return c.publisher.Publish(context.Background(), failEvent.EventName(), failEvent)
	}

	successEvent := domain.InventoryReservedEvent{
		OrderID:   ev.OrderID,
		CreatedAt: time.Now().UTC(),
	}
	return c.publisher.Publish(context.Background(), successEvent.EventName(), successEvent)
}
