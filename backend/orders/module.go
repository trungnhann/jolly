package orders

import (
	"context"
	"embed"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	ordersdb "jolly/backend/orders/adapters/db"
	ordersqueue "jolly/backend/orders/adapters/queue"
	ordershttp "jolly/backend/orders/api/http"
	ordersmodule "jolly/backend/orders/api/module"
	"jolly/backend/orders/app/command"
	"jolly/backend/orders/app/query"
)

type Module struct {
	pgxDb *pgxpool.Pool

	contracts       *contracts.Contracts
	publisher       message.Publisher
	subscriber      message.Subscriber
	commandHandlers *command.Handlers
	queryHandlers   *query.Handlers
}

func NewModule(pgxDb *pgxpool.Pool, contracts *contracts.Contracts, publisher message.Publisher, subscriber message.Subscriber) *Module {
	return &Module{
		pgxDb:      pgxDb,
		contracts:  contracts,
		publisher:  publisher,
		subscriber: subscriber,
	}
}

func (m *Module) Name() module.Name {
	return "orders"
}

//go:embed adapters/db/migrations/*.sql
var embedMigrations embed.FS

func (m *Module) Init(ctx context.Context) error {
	if err := common.MigrateDatabaseUp(
		ctx,
		string(m.Name()),
		m.pgxDb,
		embedMigrations,
		"adapters/db/migrations",
	); err != nil {
		return err
	}

	orderRepository := ordersdb.NewPostgresRepository(m.pgxDb)
	publisher := ordersqueue.NewPublisher(m.publisher)
	m.commandHandlers = command.NewHandlers(m.contracts, orderRepository, publisher)
	m.queryHandlers = query.NewHandlers(orderRepository)
	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, contracts *contracts.Contracts) error {
	contracts.Orders = ordersmodule.New(m.commandHandlers, m.queryHandlers)
	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, e common.EchoRouter) error {
	return ordershttp.Register(ctx, e, m.commandHandlers, m.queryHandlers)
}

func (m *Module) RegisterEventHandlers(ctx context.Context, router *message.Router) error {
	consumer := ordersqueue.NewConsumer(m.commandHandlers)

	router.AddConsumerHandler(
		"orders.on_inventory_reserved",
		"inventory.reserved",
		m.subscriber,
		consumer.HandleInventoryReserved,
	)

	router.AddConsumerHandler(
		"orders.on_inventory_reservation_failed",
		"inventory.reservation.failed",
		m.subscriber,
		consumer.HandleInventoryReservationFailed,
	)

	return nil
}
