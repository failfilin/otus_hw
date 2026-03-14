package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Dish struct {
	Id       uuid.UUID `json:"id" swaggertype:"string" format:"uuid"`
	MenuID   uuid.UUID `json:"menuId" swaggertype:"string" format:"uuid"`
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

func (d *Dish) Update(name, compound *string, price *float32, category *string, macros *Macros) {
	if name != nil {
		d.Name = *name
	}
	if compound != nil {
		d.Compound = *compound
	}
	if price != nil {
		d.Price = *price
	}
	if category != nil {
		d.Category = *category
	}
	if macros != nil {
		d.Macros = *macros
	}
}
