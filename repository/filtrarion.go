package repository

import (
	"otus/models"
)

var RestSlice models.SafeSlice[models.Restaurant]
var MenuSlice models.SafeSlice[models.Menu]
var DishSlice models.SafeSlice[models.Dish]

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
