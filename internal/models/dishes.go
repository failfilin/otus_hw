package models

import "github.com/google/uuid"

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
