package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/failfilin/otus_hw/internal/repository"
	"github.com/failfilin/otus_hw/internal/server"
	"github.com/failfilin/otus_hw/internal/service"

	_ "github.com/failfilin/otus_hw/docs"

	"github.com/failfilin/otus_hw/internal/models"
)

// @title           Restaurants API
// @version         1.0
// @description     Simple REST API for restaurants
// @BasePath        /api
func main() {
	var count int
	var wgProduce, wgConsume, wgLogger sync.WaitGroup
	repository.InitRepository()
	channel := make(chan models.EatType)
	doneChannel := make(chan struct{})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	srv := server.New() // Создаем новый сервер
	if err := srv.Run(ctx, ":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	defer stop()
	fmt.Println("Введи количество итераций")
	fmt.Fscan(os.Stdin, &count)
	wgLogger.Add(1)
	go service.Logger(ctx, doneChannel, &wgLogger, repository.RestSlice.Length(), repository.DishSlice.Length(), repository.MenuSlice.Length())
	wgConsume.Add(1)
	go service.NewEventConsumer(ctx, channel, doneChannel, &wgConsume)
	for i := count; i > 0; i-- {
		select {
		case <-ctx.Done():
		default:
			wgProduce.Add(1)
			go service.GenerateModels(ctx, channel, &wgProduce)
		}
		if ctx.Err() != nil {
			break
		}

	}
	wgProduce.Wait()
	close(channel)
	wgConsume.Wait()
	wgLogger.Wait()
}
