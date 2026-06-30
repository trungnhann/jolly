package inventory

import (
	"context"
	"embed"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	inventorydb "jolly/backend/inventory/adapters/db"
	inventoryqueue "jolly/backend/inventory/adapters/queue"
	inventoryhttp "jolly/backend/inventory/api/http"
	inventorymodule "jolly/backend/inventory/api/module"
	inventorycommand "jolly/backend/inventory/app/command"
	inventoryquery "jolly/backend/inventory/app/query"
)

type Module struct {
	pgxDb           *pgxpool.Pool
	contracts       *contracts.Contracts
	commandHandlers *inventorycommand.Handlers
	queryHandlers   *inventoryquery.Handlers
	publisher       message.Publisher
	subscriber      message.Subscriber
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
	return "inventory"
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

	repo := inventorydb.NewPostgresRepository(m.pgxDb)
	m.commandHandlers = inventorycommand.NewHandlers(repo, m.contracts)
	m.queryHandlers = inventoryquery.NewHandlers(repo)
	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, contracts *contracts.Contracts) error {
	contracts.Inventory = inventorymodule.New(m.commandHandlers)
	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, e common.EchoRouter) error {
	return inventoryhttp.Register(ctx, e, m.commandHandlers, m.queryHandlers)
}

func (m *Module) RegisterEventHandlers(ctx context.Context, router *message.Router) error {
	publisher := inventoryqueue.NewPublisher(m.publisher)
	consumer := inventoryqueue.NewConsumer(m.commandHandlers, publisher)

	router.AddConsumerHandler(
		"inventory.reserve_on_order_paid",
		"order.paid",
		m.subscriber,
		consumer.HandleOrderPaid,
	)
	return nil
}
