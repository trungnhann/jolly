package queue

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"

	"jolly/backend/orders/app/command"
)

type IncomingInventoryReservedEvent struct {
	OrderID string `json:"order_id"`
}

type IncomingInventoryReservationFailedEvent struct {
	OrderID string `json:"order_id"`
	Reason  string `json:"reason"`
}

type Consumer struct {
	commandHandlers *command.Handlers
}

func NewConsumer(commandHandlers *command.Handlers) *Consumer {
	return &Consumer{commandHandlers: commandHandlers}
}

func (c *Consumer) HandleInventoryReserved(msg *message.Message) error {
	var event IncomingInventoryReservedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	err := c.commandHandlers.MarkOrderInventoryReserved(context.Background(), event.OrderID)
	if err != nil {
		return err
	}

	return c.commandHandlers.MarkOrderPaymentAuthorized(context.Background(), event.OrderID)
}

func (c *Consumer) HandleInventoryReservationFailed(msg *message.Message) error {
	var event IncomingInventoryReservationFailedEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	return c.commandHandlers.MarkOrderFailed(context.Background(), event.OrderID)
}
