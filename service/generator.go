package service

import (
	"math/rand"
	"otus/models"

	"github.com/google/uuid"
)

func GenerateModels() models.EatType {
	switch i := rand.Intn(3); i {
	case 0:
		return models.Dish{Id: rand.Intn(9999)}
	case 1:
		return models.Menu{Id: uuid.New()}
	default:
		return models.Restaurant{Id: uuid.New()}
	}
}
