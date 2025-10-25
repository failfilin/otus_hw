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

func (menu *Menu) ShowDetais() {
	fmt.Println("ID меню", menu.Id)
	fmt.Println("Id ресторана", menu.RestaurantId)
	fmt.Println("Описание", menu.Description)
	fmt.Println("Список блюд", menu.DishList)
	fmt.Println("Статус активности меню", menu.Active)
}
