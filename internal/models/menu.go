package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Menu struct {
	Id           uuid.UUID `json:"id"`
	RestaurantId uuid.UUID `json:"restaurantId"`
	Description  string    `json:"description"`
	Active       bool      `json:"active"`
	DishList     []Dish    `json:"dishList"`
}

func (menu Menu) String() string {
	//return fmt.Sprintf("ID меню: %v\nId ресторана: %v\nОписание: %v\nСписок блюд: %v\nСтатус активности меню:%t", menu.Id, menu.RestaurantId, menu.Description, menu.DishList, menu.Active)
	return fmt.Sprintf("%v", menu.Id)
}
