package service

import (
	"otus/models"
	"otus/repository"
	"sync"
)

func NewEventConsumer(EatTypeChannel <-chan models.EatType, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range EatTypeChannel {
		repository.AddSlice(i)
	}
	close(done)
}
