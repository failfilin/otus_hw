package repository

import (
	"github.com/failfilin/otus_hw/internal/models"
)

var RestSlice = models.SafeSlice[models.Restaurant]{}
var MenuSlice = models.SafeSlice[models.Menu]{}
var DishSlice = models.SafeSlice[models.Dish]{}

func AddSlice(i models.EatType) {
	switch v := i.(type) {
	case models.Restaurant:
		RestSlice.Append(v)
		InsertDocument(CollectionRestaurants, v)
	case models.Menu:
		MenuSlice.Append(v)
		InsertDocument(CollectionMenus, v)
	case models.Dish:
		DishSlice.Append(v)
		InsertDocument(CollectionDishes, v)
	}
}
