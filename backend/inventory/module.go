package inventory

import (
	"context"

	"jolly/backend/common"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	inventorymodule "jolly/backend/inventory/api/module"
	"jolly/backend/inventory/app"
)

type Module struct {
	service *app.Service
}

func NewModule() *Module {
	return &Module{}
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
