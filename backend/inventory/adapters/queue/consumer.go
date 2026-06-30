package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"jolly/backend/common/event"
	"jolly/backend/inventory/api/module/client"
	"jolly/backend/inventory/app"
	"jolly/backend/inventory/domain"
)

type IncomingOrderCreatedEvent struct {
	OrderID string `json:"order_id"`
	Items   []struct {
		SKU      string `json:"sku"`
		Quantity int    `json:"quantity"`
	} `json:"items"`
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

	req := client.ReserveStockRequest{
		OrderID: ev.OrderID,
		Items:   make([]client.ReserveStockItem, 0, len(ev.Items)),
	}
	for _, item := range ev.Items {
		req.Items = append(req.Items, client.ReserveStockItem{
			SKU:      item.SKU,
			Quantity: item.Quantity,
		})
	}

	err := c.service.Reserve(context.Background(), req)
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
