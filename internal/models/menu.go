package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	Id           uuid.UUID `json:"id" swaggertype:"string" format:"uuid"`
	RestaurantId uuid.UUID `json:"restaurantId" swaggertype:"string" format:"uuid"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Active       bool      `json:"active"`
	DishList     []Dish    `json:"dishList"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (menu Menu) String() string {
	//return fmt.Sprintf("ID меню: %v\nId ресторана: %v\nОписание: %v\nСписок блюд: %v\nСтатус активности меню:%t", menu.Id, menu.RestaurantId, menu.Description, menu.DishList, menu.Active)
	return fmt.Sprintf("%v", menu.Id)
}

func (r Menu) GetID() any {
	return r.Id
}
