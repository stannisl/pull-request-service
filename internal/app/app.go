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

	"github.com/jackc/pgx/v5"
	"github.com/stannisl/avito-test/internal/config"
	"github.com/stannisl/avito-test/internal/service"
	"github.com/stannisl/avito-test/internal/transport/http/router"
)

type App struct {
	router router.Router

	server *http.Server
	conn   *pgx.Conn
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

	// connect to DB
	conn, err := pgx.Connect(ctx, a.Config.Database.DSN)
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
	}
	a.conn = conn

	// TODO create repositories
	//repositories := repository.Dependencies{
	//
	//}

	services := service.Dependencies{
		TeamService:        service.NewTeamService(),
		UserService:        service.NewUserService(),
		PullRequestService: service.NewPullRequestService(),
	}

	a.router = router.New(services)

	return nil
}

func (a *App) StartAndServe(ctx context.Context) error {
	a.server = &http.Server{
		Addr:    a.Config.HTTPServer.Address(),
		Handler: a.router,
	}

	go func() {
		log.Printf("HTTP server listening on %s", a.Config.HTTPServer.Address())
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	waitForShutdown() // Block main goroutine

	log.Println("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultShutdownTimeout)
	defer cancel()

	return a.server.Shutdown(ctx)
}

func waitForShutdown() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
