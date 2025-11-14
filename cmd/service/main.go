package main

import (
	"context"
	"log"

	"github.com/stannisl/avito-test/internal/app"
)

func main() {
	application := app.App{}
	log.Println("Starting application...")
	err := application.Setup(context.Background())
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	err = application.StartAndServeHTTP(context.Background())
	if err != nil {
		log.Printf("Failed to start application: %v", err)
	}
	log.Println("Graceful shutdown complete")
}
