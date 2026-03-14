package server

import (
	"encoding/json"
	"net/http"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/failfilin/otus_hw/internal/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	rest, ok := repository.RestSlice.GetByID(id)
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
// @Description  Delete restaurant by ID with all its menus and dishes
// @Tags         restaurants
// @Produce      json
// @Param        id path string true "Restaurant ID"
// @Success      200 {object} map[string]string "Deleted successfully"
// @Failure      400 {string} string "Invalid ID"
// @Failure      404 {string} string "Restaurant not found"
// @Router       /restaurants/{id} [delete]
// @Security     BearerAuth
func (s *Server) handleDeleteByIDRestaurants(w http.ResponseWriter, r *http.Request) {
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

	ok := repository.RestSlice.RemoveByID(id)
	if !ok {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	// Collect menu IDs for cascade delete
	var menuIds []uuid.UUID
	repository.MenuSlice.Mu.Lock()
	for _, m := range repository.MenuSlice.Slice {
		if m.RestaurantId == id {
			menuIds = append(menuIds, m.Id)
		}
	}
	repository.MenuSlice.Mu.Unlock()

	// Delete dishes for each menu
	for _, menuId := range menuIds {
		mid := menuId
		repository.DishSlice.RemoveWhere(func(d models.Dish) bool {
			return d.MenuID == mid
		})
		repository.DeleteManyDocuments(repository.CollectionDishes, bson.D{{Key: "menuid", Value: mid.String()}})
	}

	// Delete all menus of this restaurant
	repository.MenuSlice.RemoveWhere(func(m models.Menu) bool {
		return m.RestaurantId == id
	})
	repository.DeleteManyDocuments(repository.CollectionMenus, bson.D{{Key: "restaurantid", Value: id.String()}})

	// Delete restaurant
	repository.DeleteDocument(repository.CollectionRestaurants, bson.D{{Key: "id", Value: id.String()}})

	writeJSON(w, map[string]string{"message": "All Done"})
}

// handleChangeByIdRestaurants godoc
// @Summary      Update restaurant
// @Description  Update restaurant fields by ID
// @Tags         restaurants
// @Accept       json
// @Produce      json
// @Param        id path string true "Restaurant ID"
// @Param        restaurant body models.UpdateRestaurantRequest true "Updated restaurant fields"
// @Success      200 {object} map[string]string "Restaurant updated"
// @Failure      400 {string} string "Invalid ID or JSON"
// @Failure      404 {string} string "Restaurant not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /restaurants/{id} [put]
// @Security     BearerAuth
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

	var req models.UpdateRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	rest, ok := repository.RestSlice.GetByID(id)
	if !ok {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	rest.Update(&req.Name, &req.Logo, nil, &req.Active)

	if err := repository.UpdateDocument(
		repository.CollectionRestaurants,
		bson.D{{Key: "id", Value: id.String()}},
		bson.D{{Key: "$set", Value: rest}},
	); err != nil {
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
// @Param        restaurant body models.CreateRestaurantRequest true "Restaurant payload"
// @Success      201 {object} map[string]string "Generated restaurant ID"
// @Failure      400 {string} string "Invalid JSON or missing fields"
// @Failure      500 {string} string "Internal server error"
// @Router       /restaurants [post]
// @Security     BearerAuth
func (s *Server) AddNewRestaurants(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRestaurantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	payload := models.Restaurant{
		Id:     uuid.New(),
		Name:   req.Name,
		Logo:   req.Logo,
		Active: false,
	}

	repository.RestSlice.Append(payload)
	if err := repository.InsertDocument(repository.CollectionRestaurants, payload); err != nil {
		http.Error(w, "failed to save to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"id": payload.Id.String()})
}

// --- Menu handlers ---

// handleCreateMenu godoc
// @Summary      Create menu
// @Description  Create new menu for a restaurant
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        menu body models.CreateMenuRequest true "Menu payload"
// @Success      201 {object} map[string]string "Generated menu ID"
// @Failure      400 {string} string "Invalid JSON"
// @Failure      404 {string} string "Restaurant not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /menus [post]
// @Security     BearerAuth
func (s *Server) handleCreateMenu(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	restaurantId, err := uuid.Parse(req.RestaurantId)
	if err != nil {
		http.Error(w, "invalid restaurantId", http.StatusBadRequest)
		return
	}

	if _, ok := repository.RestSlice.GetByID(restaurantId); !ok {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	payload := models.Menu{
		Id:           uuid.New(),
		RestaurantId: restaurantId,
		Description:  req.Description,
		Active:       req.Active,
	}

	repository.MenuSlice.Append(payload)
	if err := repository.InsertDocument(repository.CollectionMenus, payload); err != nil {
		http.Error(w, "failed to save to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"id": payload.Id.String()})
}

// handleMenusByRestaurant godoc
// @Summary      Get all menus for restaurant
// @Description  Returns all menus for a restaurant with active status (no dish list)
// @Tags         menus
// @Produce      json
// @Param        id path string true "Restaurant ID"
// @Success      200 {array} models.Menu
// @Failure      400 {string} string "Invalid ID"
// @Router       /restaurants/{id}/menus [get]
func (s *Server) handleMenusByRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	repository.MenuSlice.Mu.Lock()
	defer repository.MenuSlice.Mu.Unlock()

	result := []models.Menu{}
	for _, m := range repository.MenuSlice.Slice {
		if m.RestaurantId == id {
			m.DishList = nil
			result = append(result, m)
		}
	}
	writeJSON(w, result)
}

// handleActiveMenuByRestaurant godoc
// @Summary      Get active menu for restaurant with dishes
// @Description  Returns the active menu for a restaurant including all its dishes
// @Tags         menus
// @Produce      json
// @Param        id path string true "Restaurant ID"
// @Success      200 {object} models.Menu
// @Failure      400 {string} string "Invalid ID"
// @Failure      404 {string} string "No active menu found"
// @Router       /restaurants/{id}/menu/active [get]
func (s *Server) handleActiveMenuByRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var activeMenu models.Menu
	found := false
	repository.MenuSlice.Mu.Lock()
	for _, m := range repository.MenuSlice.Slice {
		if m.RestaurantId == id && m.Active {
			activeMenu = m
			found = true
			break
		}
	}
	repository.MenuSlice.Mu.Unlock()

	if !found {
		http.Error(w, "no active menu found", http.StatusNotFound)
		return
	}

	var dishes []models.Dish
	repository.DishSlice.Mu.Lock()
	for _, d := range repository.DishSlice.Slice {
		if d.MenuID == activeMenu.Id {
			dishes = append(dishes, d)
		}
	}
	repository.DishSlice.Mu.Unlock()

	if dishes == nil {
		dishes = []models.Dish{}
	}
	activeMenu.DishList = dishes
	writeJSON(w, activeMenu)
}

// handleUpdateMenu godoc
// @Summary      Update menu
// @Description  Update menu description and/or active status
// @Tags         menus
// @Accept       json
// @Produce      json
// @Param        id path string true "Menu ID"
// @Param        menu body models.UpdateMenuRequest true "Updated menu fields"
// @Success      200 {object} map[string]string
// @Failure      400 {string} string "Invalid ID or JSON"
// @Failure      404 {string} string "Menu not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /menus/{id} [put]
// @Security     BearerAuth
func (s *Server) handleUpdateMenu(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req models.UpdateMenuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	menu, ok := repository.MenuSlice.GetByID(id)
	if !ok {
		http.Error(w, "menu not found", http.StatusNotFound)
		return
	}

	menu.Update(&req.Description, &req.Active)

	if err := repository.UpdateDocument(
		repository.CollectionMenus,
		bson.D{{Key: "id", Value: id.String()}},
		bson.D{{Key: "$set", Value: menu}},
	); err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "menu updated successfully"})
}

// handleDeleteMenu godoc
// @Summary      Delete menu
// @Description  Delete menu by ID along with all its dishes
// @Tags         menus
// @Produce      json
// @Param        id path string true "Menu ID"
// @Success      200 {object} map[string]string
// @Failure      400 {string} string "Invalid ID"
// @Failure      404 {string} string "Menu not found"
// @Router       /menus/{id} [delete]
// @Security     BearerAuth
func (s *Server) handleDeleteMenu(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if ok := repository.MenuSlice.RemoveByID(id); !ok {
		http.Error(w, "menu not found", http.StatusNotFound)
		return
	}

	repository.DishSlice.RemoveWhere(func(d models.Dish) bool {
		return d.MenuID == id
	})
	repository.DeleteManyDocuments(repository.CollectionDishes, bson.D{{Key: "menuid", Value: id.String()}})
	repository.DeleteDocument(repository.CollectionMenus, bson.D{{Key: "id", Value: id.String()}})

	writeJSON(w, map[string]string{"message": "menu deleted"})
}

// --- Dish handlers ---

// handleCreateDish godoc
// @Summary      Create dish
// @Description  Create new dish for a menu
// @Tags         dishes
// @Accept       json
// @Produce      json
// @Param        dish body models.CreateDishRequest true "Dish payload"
// @Success      201 {object} map[string]string "Generated dish ID"
// @Failure      400 {string} string "Invalid JSON or missing fields"
// @Failure      404 {string} string "Menu not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /dishes [post]
// @Security     BearerAuth
func (s *Server) handleCreateDish(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	menuId, err := uuid.Parse(req.MenuID)
	if err != nil {
		http.Error(w, "invalid menuId", http.StatusBadRequest)
		return
	}

	if _, ok := repository.MenuSlice.GetByID(menuId); !ok {
		http.Error(w, "menu not found", http.StatusNotFound)
		return
	}

	payload := models.Dish{
		Id:       uuid.New(),
		MenuID:   menuId,
		Name:     req.Name,
		Compound: req.Compound,
		Price:    req.Price,
		Category: req.Category,
		Macros: models.Macros{
			Calories:      req.Macros.Calories,
			Proteins:      req.Macros.Proteins,
			Fats:          req.Macros.Fats,
			Carbohydrates: req.Macros.Carbohydrates,
		},
	}

	repository.DishSlice.Append(payload)
	if err := repository.InsertDocument(repository.CollectionDishes, payload); err != nil {
		http.Error(w, "failed to save to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"id": payload.Id.String()})
}

// handleDishByID godoc
// @Summary      Get dish by ID
// @Description  Returns dish by its ID
// @Tags         dishes
// @Produce      json
// @Param        id path string true "Dish ID"
// @Success      200 {object} models.Dish
// @Failure      400 {string} string "Invalid ID"
// @Failure      404 {string} string "Dish not found"
// @Router       /dishes/{id} [get]
func (s *Server) handleDishByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	dish, ok := repository.DishSlice.GetByID(id)
	if !ok {
		http.Error(w, "dish not found", http.StatusNotFound)
		return
	}
	writeJSON(w, dish)
}

// handleUpdateDish godoc
// @Summary      Update dish
// @Description  Update dish fields by ID
// @Tags         dishes
// @Accept       json
// @Produce      json
// @Param        id path string true "Dish ID"
// @Param        dish body models.UpdateDishRequest true "Updated dish fields"
// @Success      200 {object} map[string]string
// @Failure      400 {string} string "Invalid ID or JSON"
// @Failure      404 {string} string "Dish not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /dishes/{id} [put]
// @Security     BearerAuth
func (s *Server) handleUpdateDish(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req models.UpdateDishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	dish, ok := repository.DishSlice.GetByID(id)
	if !ok {
		http.Error(w, "dish not found", http.StatusNotFound)
		return
	}

	macros := models.Macros{
		Calories:      req.Macros.Calories,
		Proteins:      req.Macros.Proteins,
		Fats:          req.Macros.Fats,
		Carbohydrates: req.Macros.Carbohydrates,
	}
	dish.Update(&req.Name, &req.Compound, &req.Price, &req.Category, &macros)

	if err := repository.UpdateDocument(
		repository.CollectionDishes,
		bson.D{{Key: "id", Value: id.String()}},
		bson.D{{Key: "$set", Value: dish}},
	); err != nil {
		http.Error(w, "failed to update", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "dish updated successfully"})
}

// --- Cart handler ---

// handleSaveCart godoc
// @Summary      Save cart
// @Description  Save a cart (no auth required). Calculates total if status is closed.
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        cart body models.SaveCartRequest true "Cart payload"
// @Success      201 {object} map[string]string "Generated cart ID"
// @Failure      400 {string} string "Invalid JSON"
// @Failure      404 {string} string "Restaurant not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /cart [post]
func (s *Server) handleSaveCart(w http.ResponseWriter, r *http.Request) {
	var req models.SaveCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	restaurantId, err := uuid.Parse(req.RestaurantId)
	if err != nil {
		http.Error(w, "invalid restaurantId", http.StatusBadRequest)
		return
	}

	if _, ok := repository.RestSlice.GetByID(restaurantId); !ok {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	var dishIds []uuid.UUID
	for _, idStr := range req.DishIds {
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid dishId: "+idStr, http.StatusBadRequest)
			return
		}
		dishIds = append(dishIds, id)
	}

	status := models.CartStatus(req.Status)
	if status == "" {
		status = models.CartStatusOpen
	}

	payload := models.Cart{
		Id:           uuid.New(),
		RestaurantId: restaurantId,
		DishIds:      dishIds,
		Status:       status,
	}

	// Calculate total when closing cart
	if payload.Status == models.CartStatusClosed {
		var total float32
		for _, dishId := range payload.DishIds {
			if dish, ok := repository.DishSlice.GetByID(dishId); ok {
				total += dish.Price
			}
		}
		payload.Total = total
	}

	repository.CartSlice.Append(payload)
	if err := repository.InsertDocument(repository.CollectionCarts, payload); err != nil {
		http.Error(w, "failed to save to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"id": payload.Id.String()})
}
