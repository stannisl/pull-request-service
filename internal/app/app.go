package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stannisl/avito-test/internal/config"
	"github.com/stannisl/avito-test/internal/repository"
	"github.com/stannisl/avito-test/internal/server"
	"github.com/stannisl/avito-test/internal/service"
	"github.com/stannisl/avito-test/internal/transport/http/router"
	"github.com/stannisl/avito-test/pkg/db"
)

type App struct {
	router router.Router

	server *http.Server
	db     *sqlx.DB
	Config *config.Config
}

func (a *App) Setup(ctx context.Context) error {
	// getting config
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}
	a.Config = appConfig
	log.Println("Configuration loaded")

	pool, err := db.ConnectPoolWithRetry(
		ctx, &db.OptionsDB{
			ConnStr:       a.Config.Database.ConnStr,
			MaxRetries:    a.Config.Database.Retries,
			RetryInterval: time.Second * time.Duration(a.Config.Database.RetrySecondsDelay),
			DriverName:    a.Config.Database.DriverName,
		})
	if err != nil {
		return fmt.Errorf("error connecting to database, %s", err)
	}
	a.db = pool
	log.Println("Connection pool created")

	// Getting conn for migration
	conn, releaseConn, err := db.GetConnFromPool(ctx, a.db)
	if err != nil {
		return fmt.Errorf("error connecting to database, %s", err)
	}
	defer func() {
		if err := releaseConn(); err != nil {
			log.Printf("Error migration conn from database, %s\n", err)
		}
	}()

	migrator := db.NewMigrator(conn, releaseConn)
	// Run migrations

	// WARNING REMOVE ON PROD
	migrator.Drop(ctx)
	// WARNING

	if err := migrator.Run(ctx); err != nil {
		return fmt.Errorf("error running migrations: %s", err)
	}
	log.Println("Migrations applied")

	txManager := db.NewTransactionManager(a.db)

	repositories := repository.Dependencies{
		PullRequestRepository: repository.NewPullRequestRepository(a.db, txManager),
		TeamRepository:        repository.NewTeamRepository(a.db, txManager),
		UserRepository:        repository.NewUserRepository(a.db, txManager),
	}
	services := service.Dependencies{
		TeamService: service.NewTeamService(repositories.UserRepository, repositories.TeamRepository),
		UserService: service.NewUserService(repositories.UserRepository, repositories.PullRequestRepository),
		PullRequestService: service.NewPullRequestService(
			repositories.PullRequestRepository,
			repositories.UserRepository,
			repositories.TeamRepository,
		),
	}

	a.router = router.New(services)
	log.Println("Router created")

	return nil
}

func (a *App) StartAndServeHTTP(ctx context.Context) error {
	// closing db of connections
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database connection pool, %s\n", err)
		}
	}(a.db)

	serv, timeout, err := server.NewBuilder().
		WithHandler(a.router).
		WithHost(a.Config.HTTPServer.Host).
		WithPort(a.Config.HTTPServer.Port).
		Build()
	if err != nil {
		return fmt.Errorf("error creating http server, %s", err)
	}
	a.server = serv
	log.Println("Server build successful")

	// TODO добавить закрытие сервера при ошибке от ListenAndServe

	go func() {
		log.Printf("HTTP server listening on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	waitForShutdown() // block main goroutine

	log.Println("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return a.server.Shutdown(ctx)
}

func waitForShutdown() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
