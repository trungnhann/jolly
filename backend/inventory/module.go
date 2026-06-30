package inventory

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"jolly/backend/common"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	inventoryqueue "jolly/backend/inventory/adapters/queue"
	inventorymodule "jolly/backend/inventory/api/module"
	"jolly/backend/inventory/app"
)

type Module struct {
	service    *app.Service
	publisher  message.Publisher
	subscriber message.Subscriber
}

func NewModule(publisher message.Publisher, subscriber message.Subscriber) *Module {
	return &Module{
		publisher:  publisher,
		subscriber: subscriber,
	}
}

func (m *Module) Name() module.Name {
	return "inventory"
}

func (m *Module) Init(ctx context.Context) error {
	m.service = app.NewService()
	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, contracts *contracts.Contracts) error {
	contracts.Inventory = inventorymodule.New(m.service)
	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, e common.EchoRouter) error {
	_ = e
	return nil
}

func (m *Module) RegisterEventHandlers(ctx context.Context, router *message.Router) error {
	publisher := inventoryqueue.NewPublisher(m.publisher)
	consumer := inventoryqueue.NewConsumer(m.service, publisher)

	router.AddConsumerHandler(
		"inventory.reserve_on_order_created",
		"order.created",
		m.subscriber,
		consumer.HandleOrderCreated,
	)
	return nil
}
