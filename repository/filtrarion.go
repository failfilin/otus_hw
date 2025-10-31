package repository

import (
	"otus/models"
)

var RestSlice []models.Restaurant
var MenuSlice []models.Menu
var DishSlice []models.Dish

func AddSlice(i models.EatType) {
	switch v := i.(type) {
	case models.Restaurant:
		RestSlice = append(RestSlice, v)
	case models.Menu:
		MenuSlice = append(MenuSlice, v)
	case models.Dish:
		DishSlice = append(DishSlice, v)
	}

}
