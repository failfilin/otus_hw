package service

import (
	"context"
	"otus/models"
	"otus/repository"
	"sync"
)

func NewEventConsumer(ctx context.Context, EatTypeChannel <-chan models.EatType, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return
	default:
		for i := range EatTypeChannel {
			repository.AddSlice(i)
		}
		close(done)
	}
}
