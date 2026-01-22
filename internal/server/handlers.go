package server

import (
	"encoding/json"
	"net/http"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/failfilin/otus_hw/internal/repository"
	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// handleAllRestaurants godoc
// @Summary      Get all restaurants
// @Description  Returns list of all restaurants
// @Tags         restaurants
// @Produce      json
// @Success      200 {array} models.Restaurant
// @Router       /restaurants [get]
func (s *Server) handleAllRestaurants(w http.ResponseWriter, r *http.Request) {
	repository.RestSlice.Mu.Lock()
	defer repository.RestSlice.Mu.Unlock()
	writeJSON(w, repository.RestSlice.Slice)
}

// handleByIdRestaurant godoc
// @Summary      Get restaurant by ID
// @Description  Returns restaurant by its ID
// @Tags         restaurants
// @Produce      json
// @Param        id path string true "Restaurant ID"
// @Success      200 {object} models.Restaurant
// @Failure      400 {string} string "Invalid ID"
// @Failure      404 {string} string "Restaurant not found"
// @Router       /restaurants/{id} [get]
func (s *Server) handleByIdRestaurant(w http.ResponseWriter, r *http.Request) {

	var rest *models.Restaurant
	var ok bool
	idStr := r.PathValue("id")

	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	rest, ok = repository.RestSlice.GetByID(id)
	if !ok {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}
	writeJSON(w, rest)
}

// handleMenus godoc
// @Summary      Get all menus
// @Description  Returns list of all menus
// @Tags         menus
// @Produce      json
// @Success      200 {array} models.Menu
// @Router       /menus [get]
func (s *Server) handleMenus(w http.ResponseWriter, r *http.Request) {
	repository.MenuSlice.Mu.Lock()
	defer repository.MenuSlice.Mu.Unlock()
	writeJSON(w, repository.MenuSlice.Slice)
}

// handleDishes godoc
// @Summary      Get all dishes
// @Description  Returns list of all dishes
// @Tags         dishes
// @Produce      json
// @Success      200 {array} models.Dish
// @Router       /dishes [get]
func (s *Server) handleDishes(w http.ResponseWriter, r *http.Request) {
	repository.DishSlice.Mu.Lock()
	defer repository.DishSlice.Mu.Unlock()
	writeJSON(w, repository.DishSlice.Slice)
}

// handleDeleteByIDRestaurants godoc
// @Summary      Delete restaurant
// @Description  Delete restaurant by ID
// @Tags         restaurants
// @Produce      json
// @Param        id path string true "Restaurant ID"
// @Success      200 {object} map[string]string "Deleted successfully"
// @Failure      400 {string} string "Invalid ID"
// @Failure      404 {string} string "Restaurant not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /restaurants/{id} [delete]
func (s *Server) handleDeleteByIDRestaurants(w http.ResponseWriter, r *http.Request) {
	var ok bool
	idStr := r.PathValue("id")

	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	ok = repository.RestSlice.RemoveByID(id)
	if !ok {
		http.Error(w, "something bad", http.StatusNotFound)
		return
	}
	repository.SaveToFile(&repository.RestSlice)
	writeJSON(w, map[string]string{
		"message": "All Done",
	})
}

// handleChangeByIdRestaurants godoc
// @Summary      Update restaurant
// @Description  Update restaurant fields by ID
// @Tags         restaurants
// @Accept       json
// @Produce      json
// @Param        id path string true "Restaurant ID"
// @Param        restaurant body models.Restaurant true "Updated restaurant fields"
// @Success      200 {object} map[string]string "Restaurant updated"
// @Failure      400 {string} string "Invalid ID or JSON"
// @Failure      404 {string} string "Restaurant not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /restaurants/{id} [put]
func (s *Server) handleChangeByIdRestaurants(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var payload models.Restaurant
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	rest, ok := repository.RestSlice.GetByID(id)
	if !ok {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	rest.Update(
		&payload.Name,
		&payload.Logo,
		&payload.MenuList,
		&payload.Active,
	)

	if err := repository.SaveToFile(&repository.RestSlice); err != nil {
		http.Error(w, "Oops not work, try again", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "restaurant updated successfully"})
}

// AddNewRestaurants godoc
// @Summary      Create restaurant
// @Description  Create new restaurant and return generated ID
// @Tags         restaurants
// @Accept       json
// @Produce      json
// @Param        restaurant body models.Restaurant true "Restaurant payload"
// @Success      201 {object} map[string]string "Generated restaurant ID"
// @Failure      400 {string} string "Invalid JSON or missing fields"
// @Failure      500 {string} string "Internal server error"
// @Router       /restaurants [post]
func (s *Server) AddNewRestaurants(w http.ResponseWriter, r *http.Request) {

	var payload models.Restaurant
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	payload.Id = uuid.New()

	payload.Active = false

	repository.RestSlice.Append(payload)

	// Сохраняем изменения в файл
	if err := repository.SaveToFile(&repository.RestSlice); err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	// Возвращаем Id созданного ресторан
	writeJSON(w, map[string]string{
		"id": payload.Id.String(),
	})
}
