package service

import (
	"math/rand"
	"otus/models"
	"sync"

	"github.com/google/uuid"
)

func GenerateModels(EatTypeChannel chan models.EatType, wg *sync.WaitGroup) {
	defer wg.Done()
	switch i := rand.Intn(3); i {
	case 0:
		EatTypeChannel <- models.Dish{Id: rand.Intn(9999)}
	case 1:
		EatTypeChannel <- models.Menu{Id: uuid.New()}
	default:
		EatTypeChannel <- models.Restaurant{Id: uuid.New()}
	}
}
