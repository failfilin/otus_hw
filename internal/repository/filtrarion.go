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

func AddSlice(i models.EatType) {
	switch v := i.(type) {
	case models.Restaurant:
		RestSlice.Append(v)
		SaveToFile(&RestSlice)
	case models.Menu:
		MenuSlice.Append(v)
		SaveToFile(&MenuSlice)
	case models.Dish:
		DishSlice.Append(v)
		SaveToFile(&DishSlice)
	}

}
