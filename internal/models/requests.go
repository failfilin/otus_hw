package models

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" example:"kek"`
	Password string `json:"password" example:"kek"`
}

// CreateRestaurantRequest represents payload for creating a restaurant
type CreateRestaurantRequest struct {
	Name string `json:"name" example:"Пицца у Тони"`
	Logo string `json:"logo" example:"https://cdn.example.com/logo.png"`
}

// UpdateRestaurantRequest represents payload for updating a restaurant
type UpdateRestaurantRequest struct {
	Name   string `json:"name" example:"Пицца у Тони"`
	Logo   string `json:"logo" example:"https://cdn.example.com/logo.png"`
	Active bool   `json:"active" example:"true"`
}

// CreateMenuRequest represents payload for creating a menu
type CreateMenuRequest struct {
	RestaurantId string `json:"restaurantId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Description  string `json:"description" example:"Летнее меню"`
	Active       bool   `json:"active" example:"false"`
}

// UpdateMenuRequest represents payload for updating a menu
type UpdateMenuRequest struct {
	Description string `json:"description" example:"Зимнее меню"`
	Active      bool   `json:"active" example:"true"`
}

// MacrosRequest represents macronutrients in a dish
type MacrosRequest struct {
	Calories      float32 `json:"calories" example:"350.5"`
	Proteins      float32 `json:"proteins" example:"12.3"`
	Fats          float32 `json:"fats" example:"8.7"`
	Carbohydrates float32 `json:"carbohydrates" example:"45.0"`
}

// CreateDishRequest represents payload for creating a dish
type CreateDishRequest struct {
	MenuID   string        `json:"menuId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name     string        `json:"name" example:"Маргарита"`
	Compound string        `json:"compound" example:"тесто, томат, моцарелла"`
	Price    float32       `json:"price" example:"590.00"`
	Category string        `json:"category" example:"Пицца"`
	Macros   MacrosRequest `json:"macros"`
}

// UpdateDishRequest represents payload for updating a dish
type UpdateDishRequest struct {
	Name     string        `json:"name" example:"Маргарита Делюкс"`
	Compound string        `json:"compound" example:"тесто, томат, моцарелла, базилик"`
	Price    float32       `json:"price" example:"690.00"`
	Category string        `json:"category" example:"Пицца"`
	Macros   MacrosRequest `json:"macros"`
}

// SaveCartRequest represents payload for saving a cart
type SaveCartRequest struct {
	RestaurantId string   `json:"restaurantId" example:"550e8400-e29b-41d4-a716-446655440000"`
	DishIds      []string `json:"dishIds" example:"550e8400-e29b-41d4-a716-446655440001,550e8400-e29b-41d4-a716-446655440002"`
	Status       string   `json:"status" example:"open" enums:"open,closed"`
}
