package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Menu struct {
	Id           uuid.UUID
	RestaurantId uuid.UUID
	Description  string
	Active       bool
	DishList     []Dish
}

func (menu Menu) String() string {
	return fmt.Sprintf("ID меню: %v\nId ресторана: %v\nОписание: %v\nСписок блюд: %v\nСтатус активности меню:%t", menu.Id, menu.RestaurantId, menu.Description, menu.DishList, menu.Active)
}

func (menu Menu) GenerateNewOne() EatType {
	return Menu{Id: uuid.New()}
}
