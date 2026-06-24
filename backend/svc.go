package backend

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	commonhttp "jolly/backend/common/http"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	"jolly/backend/inventory"
	"jolly/backend/orders"
	"jolly/backend/payments"
	"jolly/backend/users"
)

type Config struct {
	HTTPAddr string
}

type App struct {
	server *http.Server
	dbPgx  *pgxpool.Pool
}

func New(ctx context.Context, dbPgx *pgxpool.Pool, cfg Config) (*App, error) {
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if dbPgx == nil {
		panic("db connection pool cannot be nil")
	}

	moduleContracts := &contracts.Contracts{}
	registry := module.NewRegistry(moduleContracts)

	registry.Add(
		payments.NewModule(),
		inventory.NewModule(),
		users.NewModule(dbPgx),
		orders.NewModule(dbPgx, moduleContracts),
	)

	e := commonhttp.NewEcho()

	if err := registry.InitAll(ctx); err != nil {
		return nil, err
	}
	if err := registry.RegisterContractsAll(ctx); err != nil {
		return nil, err
	}
	if err := moduleContracts.Verify(); err != nil {
		return nil, err
	}
	if err := registry.RegisterHttpAll(ctx, e); err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           e,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		server: server,
		dbPgx:  dbPgx,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	defer a.dbPgx.Close()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = a.server.Shutdown(shutdownCtx)
	}()

	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
