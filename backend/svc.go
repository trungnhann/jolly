package backend

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"jolly/backend/common/file"
	commonhttp "jolly/backend/common/http"
	"jolly/backend/common/messaging"
	"jolly/backend/common/module"
	"jolly/backend/common/module/contracts"
	"jolly/backend/inventory"
	"jolly/backend/orders"
	"jolly/backend/payments"
	"jolly/backend/products"
	"jolly/backend/users"
)

type Config struct {
	HTTPAddr string
}

type App struct {
	server *http.Server
	dbPgx  *pgxpool.Pool
	router *message.Router
}

func New(ctx context.Context, dbPgx *pgxpool.Pool, cfg Config) (*App, error) {
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if dbPgx == nil {
		panic("db connection pool cannot be nil")
	}

	storageRootDir := os.Getenv("STORAGE_ROOT_DIR")
	if storageRootDir == "" {
		storageRootDir = "/tmp/jolly_storage"
	}
	storageBaseURL := os.Getenv("STORAGE_BASE_URL")
	if storageBaseURL == "" {
		storageBaseURL = "http://localhost:8080/public"
	}
	fileStorage := file.NewPublicStorage(storageRootDir, storageBaseURL)

	moduleContracts := &contracts.Contracts{}

	logger := watermill.NewStdLogger(false, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}
	router.AddMiddleware(middleware.Recoverer)

	registry := module.NewRegistry(moduleContracts, router)

	publisher, err := messaging.NewRabbitMQPublisher(logger)
	if err != nil {
		return nil, err
	}
	subscriber, err := messaging.NewRabbitMQSubscriber(logger)
	if err != nil {
		return nil, err
	}

	registry.Add(
		payments.NewModule(),
		inventory.NewModule(publisher, subscriber),
		products.NewModule(dbPgx, moduleContracts, fileStorage),
		users.NewModule(dbPgx, fileStorage),
		orders.NewModule(dbPgx, moduleContracts, publisher, subscriber),
	)

	e := commonhttp.NewEcho()
	e.Static("/public", storageRootDir)

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
	if err := registry.RegisterEventHandlersAll(ctx); err != nil {
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
		router: router,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	defer a.dbPgx.Close()

	go func() {
		if err := a.router.Run(ctx); err != nil {
			panic(err)
		}
	}()
	<-a.router.Running()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = a.server.Shutdown(shutdownCtx)
		_ = a.router.Close()
	}()

	err := a.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
