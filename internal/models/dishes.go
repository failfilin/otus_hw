package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Dish struct {
	Id       int       `json:"id"`
	MenuID   uuid.UUID `json:"menuId"`
	Name     string    `json:"name"`
	Compound string    `json:"compound"`
	Macros   Macros    `json:"macros"`
	Price    float32   `json:"price"`
	Category string    `json:"category"`
}

type Macros struct {
	Calories      float32 `json:"calories"`
	Proteins      float32 `json:"proteins"`
	Fats          float32 `json:"fats"`
	Carbohydrates float32 `json:"carbohydrates"`
}

func (dish Dish) String() string {
	return fmt.Sprintf("%v", dish.Id)
}
func (r Dish) GetID() any {
	return r.Id
}
