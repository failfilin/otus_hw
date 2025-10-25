package main

import (
	"otus/internal/models"

	"github.com/google/uuid"
)

func main() {
	rest := models.Restaurant{
		Id:   uuid.New(),
		Name: "Веселые истории",
		Logo: "Когда нить будет ссылка"}
	firstMenu := models.Menu{
		Id:           uuid.New(),
		RestaurantId: rest.Id,
		Description:  "Первая тестовая менюха",
		Active:       true}
	firstDish := models.Dish{
		Id:       1,
		MenuID:   firstMenu.Id,
		Name:     "Веселая похлебка",
		Compound: "Остатки со столоа 1 шт",
		Macros:   models.Macros{Calories: 150.0, Proteins: 2.0, Fats: 4.0, Carbohydrates: 5.0},
		Price:    200.0}

	rest.MenuList = append(rest.MenuList, firstMenu)
	firstMenu.DishList = append(firstMenu.DishList, firstDish)
	rest.ShowDetais()
	rest.ChangeActive()
	rest.ShowDetais()
	firstMenu.ShowDetais()
}
