package repository

import (
	"otus/models"
)

func AddSlice(i models.EatType, dishSlice *[]models.EatType, menuSlice *[]models.EatType, restSlice *[]models.EatType) {
	switch i.(type) {
	case models.Restaurant:
		*restSlice = append(*restSlice, i)
	case models.Menu:
		*menuSlice = append(*menuSlice, i)
	case models.Dish:
		*dishSlice = append(*dishSlice, i)
	}

}
