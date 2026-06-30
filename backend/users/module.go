package users

import (
"github.com/ThreeDotsLabs/watermill/message"
	"context"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common"
	"jolly/backend/common/file"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	usersdb "jolly/backend/users/adapters/db"
	usershttp "jolly/backend/users/api/http"
	usersmodule "jolly/backend/users/api/module"
	"jolly/backend/users/app/command"
	"jolly/backend/users/app/query"
)

type Module struct {
	pgxDb   *pgxpool.Pool
	storage file.Storage

	commandHandlers *command.Handlers
	queryHandlers   *query.Handlers
}

func NewModule(pgxDb *pgxpool.Pool, storage file.Storage) *Module {
	return &Module{
		pgxDb:   pgxDb,
		storage: storage,
	}
}

func (m *Module) Name() module.Name {
	return "users"
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

	userRepository := usersdb.NewPostgresRepository(m.pgxDb)
	m.commandHandlers = command.NewHandlers(userRepository)
	m.queryHandlers = query.NewHandlers(userRepository)
	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, contracts *contracts.Contracts) error {
	_ = ctx
	contracts.Users = usersmodule.New(m.commandHandlers, m.queryHandlers)
	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, e common.EchoRouter) error {
	return usershttp.Register(ctx, e, m.commandHandlers, m.queryHandlers, m.storage)
}

func (m *Module) RegisterEventHandlers(ctx context.Context, router *message.Router) error {
	return nil
}
