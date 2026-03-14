package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) routes() {
	apiMux := http.NewServeMux()

	// Auth
	apiMux.HandleFunc("POST /auth/login", s.handleLogin)

	// Restaurants (mutations require JWT)
	apiMux.Handle("POST /restaurants", jwtMiddleware(http.HandlerFunc(s.AddNewRestaurants)))
	apiMux.HandleFunc("GET /restaurants", s.handleAllRestaurants)
	apiMux.HandleFunc("GET /restaurants/{id}", s.handleByIdRestaurant)
	apiMux.Handle("PUT /restaurants/{id}", jwtMiddleware(http.HandlerFunc(s.handleChangeByIdRestaurants)))
	apiMux.Handle("DELETE /restaurants/{id}", jwtMiddleware(http.HandlerFunc(s.handleDeleteByIDRestaurants)))

	// Menus (mutations require JWT)
	apiMux.Handle("POST /menus", jwtMiddleware(http.HandlerFunc(s.handleCreateMenu)))
	apiMux.HandleFunc("GET /menus", s.handleMenus)
	apiMux.HandleFunc("GET /restaurants/{id}/menus", s.handleMenusByRestaurant)
	apiMux.HandleFunc("GET /restaurants/{id}/menu/active", s.handleActiveMenuByRestaurant)
	apiMux.Handle("PUT /menus/{id}", jwtMiddleware(http.HandlerFunc(s.handleUpdateMenu)))
	apiMux.Handle("DELETE /menus/{id}", jwtMiddleware(http.HandlerFunc(s.handleDeleteMenu)))

	// Dishes (mutations require JWT)
	apiMux.Handle("POST /dishes", jwtMiddleware(http.HandlerFunc(s.handleCreateDish)))
	apiMux.HandleFunc("GET /dishes", s.handleDishes)
	apiMux.HandleFunc("GET /dishes/{id}", s.handleDishByID)
	apiMux.Handle("PUT /dishes/{id}", jwtMiddleware(http.HandlerFunc(s.handleUpdateDish)))

	// Cart (no auth required)
	apiMux.HandleFunc("POST /cart", s.handleSaveCart)

	s.mux.Handle("/api/", http.StripPrefix("/api", apiMux))
	s.mux.Handle("/swagger/", httpSwagger.WrapHandler)
}
