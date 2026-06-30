package products

import (
"github.com/ThreeDotsLabs/watermill/message"
	"context"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/common/file"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	productsdb "jolly/backend/products/adapters/db"
	productshttp "jolly/backend/products/api/http"
	productsmodule "jolly/backend/products/api/module"
	"jolly/backend/products/app/command"
	"jolly/backend/products/app/query"
)

type Module struct {
	pgxDb   *pgxpool.Pool
	storage file.Storage

	contracts       *contracts.Contracts
	commandHandlers *command.Handlers
	queryHandlers   *query.Handlers
}

func NewModule(pgxDb *pgxpool.Pool, contracts *contracts.Contracts, storage file.Storage) *Module {
	return &Module{
		pgxDb:     pgxDb,
		contracts: contracts,
		storage:   storage,
	}
}

func (m *Module) Name() module.Name {
	return "products"
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

	productRepository := productsdb.NewPostgresRepository(m.pgxDb)
	m.commandHandlers = command.NewHandlers(productRepository)
	m.queryHandlers = query.NewHandlers(productRepository)
	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, contracts *contracts.Contracts) error {
	contracts.Products = productsmodule.New(m.queryHandlers)
	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, e common.EchoRouter) error {
	return productshttp.Register(ctx, e, m.commandHandlers, m.queryHandlers, m.storage)
}

func (m *Module) RegisterEventHandlers(ctx context.Context, router *message.Router) error {
	return nil
}
