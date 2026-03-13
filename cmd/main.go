package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/failfilin/otus_hw/internal/repository"
	"github.com/failfilin/otus_hw/internal/server"

	_ "github.com/failfilin/otus_hw/docs"
)

// @title           Restaurants API
// @version         1.0
// @descriptionТы     Simple REST API for restaurants
// @BasePath        /api
func main() {
	if err := repository.InitPostgres(); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	repository.LoadCollectionsFromPostgres()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New()
	if err := srv.Run(ctx, ":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
