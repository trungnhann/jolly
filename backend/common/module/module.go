package module

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"

	"jolly/backend/common"
	"jolly/backend/common/module/contracts"
)

type Name string

type Module interface {
	Name() Name
	Init(ctx context.Context) error
	RegisterContracts(ctx context.Context, contracts *contracts.Contracts) error
	RegisterHttp(ctx context.Context, e common.EchoRouter) error
	RegisterEventHandlers(ctx context.Context, router *message.Router) error
}

type Registry struct {
	modules   []Module
	contracts *contracts.Contracts
	router    *message.Router
}

func NewRegistry(contracts *contracts.Contracts, router *message.Router) *Registry {
	return &Registry{
		modules:   []Module{},
		contracts: contracts,
		router:    router,
	}
}

func (r *Registry) Add(modules ...Module) {
	r.modules = append(r.modules, modules...)
}

func (r *Registry) InitAll(ctx context.Context) error {
	for _, m := range r.modules {
		if err := m.Init(ctx); err != nil {
			return fmt.Errorf("%s init: %w", m.Name(), err)
		}
	}
	return nil
}

func (r *Registry) RegisterContractsAll(ctx context.Context) error {
	for _, m := range r.modules {
		if err := m.RegisterContracts(ctx, r.contracts); err != nil {
			return fmt.Errorf("%s register contracts: %w", m.Name(), err)
		}
	}
	return nil
}

func (r *Registry) RegisterHttpAll(ctx context.Context, e common.EchoRouter) error {
	for _, m := range r.modules {
		if err := m.RegisterHttp(ctx, e); err != nil {
			return fmt.Errorf("%s register http: %w", m.Name(), err)
		}
	}
	return nil
}

func (r *Registry) RegisterEventHandlersAll(ctx context.Context) error {
	for _, m := range r.modules {
		if err := m.RegisterEventHandlers(ctx, r.router); err != nil {
			return fmt.Errorf("%s register event handlers: %w", m.Name(), err)
		}
	}
	return nil
}
