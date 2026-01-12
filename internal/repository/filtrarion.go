package repository

import (
	"github.com/failfilin/otus_hw/internal/models"
)

var RestSlice = models.SafeSlice[models.Restaurant]{
	File: "data/restaurant.json",
}
var MenuSlice = models.SafeSlice[models.Menu]{
	File: "data/menu.json",
}
var DishSlice = models.SafeSlice[models.Dish]{
	File: "data/dishes.json",
}

func InitRepository() {
	_ = RestSlice.LoadFromFile()
	_ = MenuSlice.LoadFromFile()
	_ = DishSlice.LoadFromFile()
}

func AddSlice(i models.EatType) {
	switch v := i.(type) {
	case models.Restaurant:
		RestSlice.Append(v)
	case models.Menu:
		MenuSlice.Append(v)
	case models.Dish:
		DishSlice.Append(v)
	}

}
