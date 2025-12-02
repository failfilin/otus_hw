package service

import (
	"context"
	"sync"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/failfilin/otus_hw/internal/repository"
)

func NewEventConsumer(ctx context.Context, EatTypeChannel <-chan models.EatType, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range EatTypeChannel {
		select {
		case <-ctx.Done():
			return
		default:
			repository.AddSlice(i)
		}
	}
	close(done)
}
