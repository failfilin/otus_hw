package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/failfilin/otus_hw/internal/repository"
	"github.com/failfilin/otus_hw/internal/server"
	"github.com/joho/godotenv"

	_ "github.com/failfilin/otus_hw/docs"
)

// @title           Restaurants API
// @version         1.0
// @description     Simple REST API for restaurants
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	repository.InitRepository()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New()
	if err := srv.Run(ctx, ":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
