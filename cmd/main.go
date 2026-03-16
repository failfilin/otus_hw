package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/failfilin/otus_hw/internal/service"
	"github.com/failfilin/otus_hw/internal/repository"
	"github.com/failfilin/otus_hw/internal/server"

	_ "github.com/failfilin/otus_hw/docs"
)

// @title           Restaurants API
// @version         1.0
// @description     Simple REST API for restaurants
// @BasePath        /api
func main() {
	repository.InitRepository()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	l := service.NewRedisLogger()
	defer l.Close()

	srv := server.New(l)
	if err := srv.Run(ctx, ":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
