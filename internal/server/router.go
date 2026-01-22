package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) routes() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("POST /restaurants", s.AddNewRestaurants)                  //создание нового ресторана
	apiMux.HandleFunc("PUT /restaurants/{id}", s.handleChangeByIdRestaurants)    //апдейт конкретного
	apiMux.HandleFunc("GET /restaurants", s.handleAllRestaurants)                //выдает все
	apiMux.HandleFunc("GET /restaurants/{id}", s.handleByIdRestaurant)           //выдает конкретный
	apiMux.HandleFunc("DELETE /restaurants/{id}", s.handleDeleteByIDRestaurants) //удаляет элемент
	apiMux.HandleFunc("GET /menus", s.handleMenus)
	apiMux.HandleFunc("GET /dishes", s.handleDishes)
	s.mux.Handle("/api/", http.StripPrefix("/api", apiMux))
	swaggerMux := http.NewServeMux()
	swaggerMux.Handle("/", httpSwagger.WrapHandler)
	s.mux.Handle("/swagger/", httpSwagger.WrapHandler)

}
