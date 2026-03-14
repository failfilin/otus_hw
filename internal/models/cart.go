package models

import "github.com/google/uuid"

type CartStatus string

const (
	CartStatusOpen   CartStatus = "open"
	CartStatusClosed CartStatus = "closed"
)

type Cart struct {
	Id           uuid.UUID   `json:"id" swaggertype:"string" format:"uuid"`
	RestaurantId uuid.UUID   `json:"restaurantId" swaggertype:"string" format:"uuid"`
	DishIds      []uuid.UUID `json:"dishIds"`
	Status       CartStatus  `json:"status"`
	Total        float32     `json:"total"`
}

func (c Cart) GetID() any {
	return c.Id
}
