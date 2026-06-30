package payments

import (
"github.com/ThreeDotsLabs/watermill/message"
	"context"

	"jolly/backend/common"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	paymentsmodule "jolly/backend/payments/api/module"
	"jolly/backend/payments/app"
)

type Module struct {
	service *app.Service
}

func NewModule() *Module {
	return &Module{}
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
	return nil
}
