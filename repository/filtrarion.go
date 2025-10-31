package repository

import (
	"otus/models"
)

func AddSlice(i models.EatType, dishSlice *[]models.EatType, menuSlice *[]models.EatType, restSlice *[]models.EatType) {
	switch i.(type) {
	case models.Restaurant:
		i.AddToSlice(restSlice)
	case models.Menu:
		i.AddToSlice(menuSlice)
	case models.Dish:
		i.AddToSlice(dishSlice)
	}

}
