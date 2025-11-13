package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mostanin/avito-test/internal/app"
	"github.com/mostanin/avito-test/internal/config"
)

const (
	defaultShutdownTimeout = 10 * time.Second
)

func main() {
	cfg := config.MustLoad()

	application, err := app.New(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address(),
		Handler:      application.Router(),
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		log.Printf("HTTP server listening on %s", cfg.HTTPServer.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	waitForShutdown()

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func waitForShutdown() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
}
