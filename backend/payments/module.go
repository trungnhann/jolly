package payments

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"jolly/backend/common"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	paymentsqueue "jolly/backend/payments/adapters/queue"
	paymentsmodule "jolly/backend/payments/api/module"
	"jolly/backend/payments/app"
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
	return "payments"
}

func (m *Module) Init(ctx context.Context) error {
	m.service = app.NewService()
	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, contracts *contracts.Contracts) error {
	contracts.Payments = paymentsmodule.New(m.service)
	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, e common.EchoRouter) error {
	_ = e
	return nil
}

func (m *Module) RegisterEventHandlers(ctx context.Context, router *message.Router) error {
	pub := paymentsqueue.NewPublisher(m.publisher)
	consumer := paymentsqueue.NewConsumer(m.service, pub)

	router.AddConsumerHandler(
		"payments.authorize_on_order_created",
		"order.created",
		m.subscriber,
		consumer.HandleOrderCreated,
	)
	return nil
}
