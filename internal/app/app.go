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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stannisl/avito-test/internal/config"
	"github.com/stannisl/avito-test/internal/server"
	"github.com/stannisl/avito-test/internal/service"
	"github.com/stannisl/avito-test/internal/transport/http/router"
	"github.com/stannisl/avito-test/pkg/db"
)

type App struct {
	router router.Router

	server *http.Server
	pool   *pgxpool.Pool
	Config *config.Config
}

func (a *App) Setup(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("context is required")
	}

	// getting config
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	a.Config = appConfig

	pool, err := db.ConnectPoolWithRetry(
		ctx, a.Config.Database.ConnStr,
		a.Config.Database.Retries,
		time.Second*time.Duration(a.Config.Database.RetrySecondsDelay))
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
	}
	a.pool = pool

	// Getting conn for migration
	conn, releaseConn, err := db.GetConnFromPool(ctx, a.pool)
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
	}
	defer releaseConn()

	// Run migrations
	if err := db.RunMigrations(ctx, conn.Conn()); err != nil {
		log.Fatalf("Error running migrations: %s", err)
	}

	migrationVersion, err := db.GetMigrationVersion(ctx, conn.Conn())
	if err != nil {
		log.Fatalf("Error getting db version: %s", err)
	}
	log.Printf("Latest Migration version: %v\n", migrationVersion)

	// Initing components
	//repositories := repository.Dependencies{
	//	PullRequestRepository: nil,
	//	TeamRepository:        nil,
	//	UserRepository:        nil,
	//}

	services := service.Dependencies{
		TeamService:        service.NewTeamService(),
		UserService:        service.NewUserService(),
		PullRequestService: service.NewPullRequestService(),
	}

	a.router = router.New(services)

	return nil
}

func (a *App) StartAndServeHTTP(ctx context.Context) error {
	// closing pool of connections
	defer a.pool.Close()

	serv, timeout, err := server.NewBuilder().
		WithHandler(a.router).
		WithHost(a.Config.HTTPServer.Host).
		WithPort(a.Config.HTTPServer.Port).
		Build()
	if err != nil {
		return err
	}
	a.server = serv

	go func() {
		log.Printf("HTTP server listening on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	waitForShutdown() // block main goroutine

	log.Println("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return a.server.Shutdown(ctx)
}

func waitForShutdown() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
