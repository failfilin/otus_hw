package repository

import (
	"otus/models"
	"sync"
)

var RestSlice models.SafeSlice[models.Restaurant]
var MenuSlice models.SafeSlice[models.Menu]
var DishSlice models.SafeSlice[models.Dish]

func AddSlice(EatTypeChannel <-chan models.EatType, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range EatTypeChannel {
		switch v := i.(type) {
		case models.Restaurant:
			RestSlice.Append(v)
		case models.Menu:
			MenuSlice.Append(v)
		case models.Dish:
			DishSlice.Append(v)
		}
	}
	close(done)
}
