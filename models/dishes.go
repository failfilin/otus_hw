package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Dish struct {
	Id       int
	MenuID   uuid.UUID
	Name     string
	Compound string
	Macros   Macros
	Price    float32
	Category string //
}

type Macros struct {
	Calories      float32
	Proteins      float32
	Fats          float32
	Carbohydrates float32
}

func (dish Dish) String() string {
	return fmt.Sprintf(dish.Name)
}

func (dish Dish) AddToSlice(slice *[]EatType) {
	*slice = append(*slice, dish)
}
