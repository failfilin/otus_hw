package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/failfilin/otus_hw/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// helpers

func newServer() *Server {
	return &Server{}
}

func seedRestaurant() models.Restaurant {
	r := models.Restaurant{
		Id:     uuid.New(),
		Name:   "Test Restaurant",
		Logo:   "logo.png",
		Active: true,
	}
	repository.RestSlice.Append(r)
	return r
}

func clearRestaurants() {
	repository.RestSlice.Mu.Lock()
	defer repository.RestSlice.Mu.Unlock()
	repository.RestSlice.Slice = nil
}

// ── handleAllRestaurants ──────────────────────────────────────────────────────

func TestHandleAllRestaurants_Empty(t *testing.T) {
	clearRestaurants()
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants", nil)
	w := httptest.NewRecorder()

	s.handleAllRestaurants(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var result []models.Restaurant
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d items", len(result))
	}
}

func TestHandleAllRestaurants_WithData(t *testing.T) {
	clearRestaurants()
	seedRestaurant()
	seedRestaurant()
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants", nil)
	w := httptest.NewRecorder()

	s.handleAllRestaurants(w, req)

	var result []models.Restaurant
	json.NewDecoder(w.Result().Body).Decode(&result)

	if len(result) != 2 {
		t.Errorf("expected 2 restaurants, got %d", len(result))
	}
}

// ── handleByIdRestaurant ──────────────────────────────────────────────────────

func TestHandleByIdRestaurant_OK(t *testing.T) {
	clearRestaurants()
	r := seedRestaurant()
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+r.Id.String(), nil)
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleByIdRestaurant(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var got models.Restaurant
	json.NewDecoder(res.Body).Decode(&got)
	if got.Id != r.Id {
		t.Errorf("expected id %s, got %s", r.Id, got.Id)
	}
}

func TestHandleByIdRestaurant_InvalidUUID(t *testing.T) {
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants/not-a-uuid", nil)
	req.SetPathValue("id", "not-a-uuid")
	w := httptest.NewRecorder()

	s.handleByIdRestaurant(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleByIdRestaurant_NotFound(t *testing.T) {
	clearRestaurants()
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+uuid.New().String(), nil)
	req.SetPathValue("id", uuid.New().String())
	w := httptest.NewRecorder()

	s.handleByIdRestaurant(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleByIdRestaurant_EmptyID(t *testing.T) {
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants/", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	s.handleByIdRestaurant(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleDeleteByIDRestaurants ───────────────────────────────────────────────

func TestHandleDeleteByIDRestaurants_OK(t *testing.T) {
	clearRestaurants()
	r := seedRestaurant()
	s := newServer()

	req := httptest.NewRequest(http.MethodDelete, "/restaurants/"+r.Id.String(), nil)
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleDeleteByIDRestaurants(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// verify actually deleted
	_, ok := repository.RestSlice.GetByID(r.Id)
	if ok {
		t.Error("restaurant should have been deleted")
	}
}

func TestHandleDeleteByIDRestaurants_InvalidUUID(t *testing.T) {
	s := newServer()

	req := httptest.NewRequest(http.MethodDelete, "/restaurants/bad-id", nil)
	req.SetPathValue("id", "bad-id")
	w := httptest.NewRecorder()

	s.handleDeleteByIDRestaurants(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleDeleteByIDRestaurants_NotFound(t *testing.T) {
	clearRestaurants()
	s := newServer()

	req := httptest.NewRequest(http.MethodDelete, "/restaurants/"+uuid.New().String(), nil)
	req.SetPathValue("id", uuid.New().String())
	w := httptest.NewRecorder()

	s.handleDeleteByIDRestaurants(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// ── handleChangeByIdRestaurants ───────────────────────────────────────────────

func TestHandleChangeByIdRestaurants_OK(t *testing.T) {
	clearRestaurants()
	r := seedRestaurant()
	s := newServer()

	updated := models.Restaurant{Name: "New Name", Logo: "new.png", Active: false}
	body, _ := json.Marshal(updated)

	req := httptest.NewRequest(http.MethodPut, "/restaurants/"+r.Id.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleChangeByIdRestaurants(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleChangeByIdRestaurants_InvalidUUID(t *testing.T) {
	s := newServer()

	req := httptest.NewRequest(http.MethodPut, "/restaurants/bad", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	s.handleChangeByIdRestaurants(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleChangeByIdRestaurants_InvalidJSON(t *testing.T) {
	clearRestaurants()
	r := seedRestaurant()
	s := newServer()

	req := httptest.NewRequest(http.MethodPut, "/restaurants/"+r.Id.String(), bytes.NewReader([]byte(`not json`)))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleChangeByIdRestaurants(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleChangeByIdRestaurants_NotFound(t *testing.T) {
	clearRestaurants()
	s := newServer()

	body, _ := json.Marshal(models.Restaurant{Name: "Ghost"})
	req := httptest.NewRequest(http.MethodPut, "/restaurants/"+uuid.New().String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", uuid.New().String())
	w := httptest.NewRecorder()

	s.handleChangeByIdRestaurants(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// ── AddNewRestaurants ─────────────────────────────────────────────────────────

func TestAddNewRestaurants_OK(t *testing.T) {
	clearRestaurants()
	s := newServer()

	payload := models.Restaurant{Name: "Brand New"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/restaurants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.AddNewRestaurants(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	json.NewDecoder(w.Result().Body).Decode(&resp)
	if resp["id"] == "" {
		t.Error("expected id in response, got empty")
	}
	if _, err := uuid.Parse(resp["id"]); err != nil {
		t.Errorf("returned id is not a valid UUID: %s", resp["id"])
	}
}

func TestAddNewRestaurants_MissingName(t *testing.T) {
	s := newServer()

	payload := models.Restaurant{Logo: "logo.png"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/restaurants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.AddNewRestaurants(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAddNewRestaurants_InvalidJSON(t *testing.T) {
	s := newServer()

	req := httptest.NewRequest(http.MethodPost, "/restaurants", bytes.NewReader([]byte(`{bad json}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.AddNewRestaurants(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAddNewRestaurants_ActiveIsAlwaysFalse(t *testing.T) {
	clearRestaurants()
	s := newServer()

	payload := models.Restaurant{Name: "Active Test", Active: true}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/restaurants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.AddNewRestaurants(w, req)

	var resp map[string]string
	json.NewDecoder(w.Result().Body).Decode(&resp)

	id, _ := uuid.Parse(resp["id"])
	saved, ok := repository.RestSlice.GetByID(id)
	if !ok {
		t.Fatal("restaurant not found after creation")
	}
	if saved.Active {
		t.Error("new restaurant should always have Active=false")
	}
}

// ── handleMenus / handleDishes ────────────────────────────────────────────────

func TestHandleMenus_ReturnsJSON(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodGet, "/menus", nil)
	w := httptest.NewRecorder()

	s.handleMenus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestHandleDishes_ReturnsJSON(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodGet, "/dishes", nil)
	w := httptest.NewRecorder()

	s.handleDishes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

// ── handleDishByID ────────────────────────────────────────────────────────────

func TestHandleDishByID_OK(t *testing.T) {
	clearDishes()
	m := seedMenu(seedRestaurant().Id)
	d := seedDish(m.Id)
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/dishes/"+d.Id.String(), nil)
	req.SetPathValue("id", d.Id.String())
	w := httptest.NewRecorder()

	s.handleDishByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var got models.Dish
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Id != d.Id {
		t.Errorf("expected id %s, got %s", d.Id, got.Id)
	}
	if got.Name != d.Name {
		t.Errorf("expected name %s, got %s", d.Name, got.Name)
	}
}

func TestHandleDishByID_NotFound_Returns404(t *testing.T) {
	clearDishes()
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/dishes/"+uuid.New().String(), nil)
	req.SetPathValue("id", uuid.New().String())
	w := httptest.NewRecorder()

	s.handleDishByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleDishByID_InvalidUUID_Returns400(t *testing.T) {
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/dishes/bad-id", nil)
	req.SetPathValue("id", "bad-id")
	w := httptest.NewRecorder()

	s.handleDishByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── Helpers (menus / dishes / carts) ─────────────────────────────────────────

func clearMenus() {
	repository.MenuSlice.Mu.Lock()
	defer repository.MenuSlice.Mu.Unlock()
	repository.MenuSlice.Slice = nil
}

func clearDishes() {
	repository.DishSlice.Mu.Lock()
	defer repository.DishSlice.Mu.Unlock()
	repository.DishSlice.Slice = nil
}

func clearCarts() {
	repository.CartSlice.Mu.Lock()
	defer repository.CartSlice.Mu.Unlock()
	repository.CartSlice.Slice = nil
}

func seedMenu(restaurantId uuid.UUID) models.Menu {
	m := models.Menu{
		Id:           uuid.New(),
		RestaurantId: restaurantId,
		Description:  "Test Menu",
		Active:       false,
	}
	repository.MenuSlice.Append(m)
	return m
}

func seedActiveMenu(restaurantId uuid.UUID) models.Menu {
	m := models.Menu{
		Id:           uuid.New(),
		RestaurantId: restaurantId,
		Description:  "Active Menu",
		Active:       true,
	}
	repository.MenuSlice.Append(m)
	return m
}

func seedDish(menuId uuid.UUID) models.Dish {
	d := models.Dish{
		Id:     uuid.New(),
		MenuID: menuId,
		Name:   "Test Dish",
		Price:  9.99,
	}
	repository.DishSlice.Append(d)
	return d
}

func validToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	str, _ := token.SignedString(jwtSecret)
	return str
}

// ── JWT middleware ─────────────────────────────────────────────────────────────

func TestJwtMiddleware_NoToken_Returns401(t *testing.T) {
	handler := jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJwtMiddleware_InvalidToken_Returns401(t *testing.T) {
	handler := jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not.a.real.token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJwtMiddleware_ValidToken_PassesThrough(t *testing.T) {
	handler := jwtMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+validToken())
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// ── handleLogin ───────────────────────────────────────────────────────────────

func TestHandleLogin_OK(t *testing.T) {
	s := newServer()
	body, _ := json.Marshal(map[string]string{"username": "kek", "password": "kek"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleLogin(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Result().Body).Decode(&resp)
	if resp["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestHandleLogin_WrongPassword_Returns401(t *testing.T) {
	s := newServer()
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	s.handleLogin(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleLogin_InvalidJSON_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(`{bad}`)))
	w := httptest.NewRecorder()

	s.handleLogin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleCreateMenu ──────────────────────────────────────────────────────────

func TestHandleCreateMenu_OK(t *testing.T) {
	clearRestaurants()
	clearMenus()
	r := seedRestaurant()
	s := newServer()

	payload := models.Menu{RestaurantId: r.Id, Description: "Lunch"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/menus", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleCreateMenu(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Result().Body).Decode(&resp)
	if _, err := uuid.Parse(resp["id"]); err != nil {
		t.Errorf("expected valid UUID id, got: %s", resp["id"])
	}
}

func TestHandleCreateMenu_RestaurantNotFound_Returns404(t *testing.T) {
	clearRestaurants()
	s := newServer()

	payload := models.Menu{RestaurantId: uuid.New(), Description: "Ghost"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/menus", bytes.NewReader(body))
	w := httptest.NewRecorder()

	s.handleCreateMenu(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleCreateMenu_InvalidJSON_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPost, "/menus", bytes.NewReader([]byte(`{bad}`)))
	w := httptest.NewRecorder()

	s.handleCreateMenu(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleMenusByRestaurant ───────────────────────────────────────────────────

func TestHandleMenusByRestaurant_OK(t *testing.T) {
	clearRestaurants()
	clearMenus()
	r := seedRestaurant()
	seedMenu(r.Id)
	seedMenu(r.Id)
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+r.Id.String()+"/menus", nil)
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleMenusByRestaurant(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var menus []models.Menu
	json.NewDecoder(w.Result().Body).Decode(&menus)
	if len(menus) != 2 {
		t.Errorf("expected 2 menus, got %d", len(menus))
	}
	for _, m := range menus {
		if m.DishList != nil {
			t.Error("dish list should be nil in this endpoint")
		}
	}
}

func TestHandleMenusByRestaurant_Empty(t *testing.T) {
	clearMenus()
	r := seedRestaurant()
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+r.Id.String()+"/menus", nil)
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleMenusByRestaurant(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var menus []models.Menu
	json.NewDecoder(w.Result().Body).Decode(&menus)
	if len(menus) != 0 {
		t.Errorf("expected 0 menus, got %d", len(menus))
	}
}

func TestHandleMenusByRestaurant_InvalidUUID_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodGet, "/restaurants/bad/menus", nil)
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	s.handleMenusByRestaurant(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleActiveMenuByRestaurant ──────────────────────────────────────────────

func TestHandleActiveMenuByRestaurant_OK(t *testing.T) {
	clearRestaurants()
	clearMenus()
	clearDishes()
	r := seedRestaurant()
	m := seedActiveMenu(r.Id)
	d := seedDish(m.Id)
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+r.Id.String()+"/menu/active", nil)
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleActiveMenuByRestaurant(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var menu models.Menu
	json.NewDecoder(w.Result().Body).Decode(&menu)
	if menu.Id != m.Id {
		t.Errorf("expected menu id %s, got %s", m.Id, menu.Id)
	}
	if len(menu.DishList) != 1 {
		t.Errorf("expected 1 dish, got %d", len(menu.DishList))
	}
	if menu.DishList[0].Id != d.Id {
		t.Errorf("expected dish id %s, got %s", d.Id, menu.DishList[0].Id)
	}
}

func TestHandleActiveMenuByRestaurant_NoActiveMenu_Returns404(t *testing.T) {
	clearRestaurants()
	clearMenus()
	r := seedRestaurant()
	seedMenu(r.Id) // inactive
	s := newServer()

	req := httptest.NewRequest(http.MethodGet, "/restaurants/"+r.Id.String()+"/menu/active", nil)
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleActiveMenuByRestaurant(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleActiveMenuByRestaurant_InvalidUUID_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodGet, "/restaurants/bad/menu/active", nil)
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	s.handleActiveMenuByRestaurant(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleUpdateMenu ──────────────────────────────────────────────────────────

func TestHandleUpdateMenu_OK(t *testing.T) {
	clearRestaurants()
	clearMenus()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	s := newServer()

	payload := models.Menu{Description: "Updated", Active: true}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/menus/"+m.Id.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", m.Id.String())
	w := httptest.NewRecorder()

	s.handleUpdateMenu(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	updated, _ := repository.MenuSlice.GetByID(m.Id)
	if updated.Description != "Updated" {
		t.Errorf("expected description 'Updated', got '%s'", updated.Description)
	}
}

func TestHandleUpdateMenu_NotFound_Returns404(t *testing.T) {
	clearMenus()
	s := newServer()

	body, _ := json.Marshal(models.Menu{Description: "x"})
	req := httptest.NewRequest(http.MethodPut, "/menus/"+uuid.New().String(), bytes.NewReader(body))
	req.SetPathValue("id", uuid.New().String())
	w := httptest.NewRecorder()

	s.handleUpdateMenu(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleUpdateMenu_InvalidUUID_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPut, "/menus/bad", bytes.NewReader([]byte(`{}`)))
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	s.handleUpdateMenu(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleUpdateMenu_InvalidJSON_Returns400(t *testing.T) {
	clearMenus()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	s := newServer()

	req := httptest.NewRequest(http.MethodPut, "/menus/"+m.Id.String(), bytes.NewReader([]byte(`{bad}`)))
	req.SetPathValue("id", m.Id.String())
	w := httptest.NewRecorder()

	s.handleUpdateMenu(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleDeleteMenu ──────────────────────────────────────────────────────────

func TestHandleDeleteMenu_OK(t *testing.T) {
	clearRestaurants()
	clearMenus()
	clearDishes()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	seedDish(m.Id)
	seedDish(m.Id)
	s := newServer()

	req := httptest.NewRequest(http.MethodDelete, "/menus/"+m.Id.String(), nil)
	req.SetPathValue("id", m.Id.String())
	w := httptest.NewRecorder()

	s.handleDeleteMenu(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if _, ok := repository.MenuSlice.GetByID(m.Id); ok {
		t.Error("menu should have been deleted")
	}
	if repository.DishSlice.Length() != 0 {
		t.Errorf("expected 0 dishes after menu delete, got %d", repository.DishSlice.Length())
	}
}

func TestHandleDeleteMenu_NotFound_Returns404(t *testing.T) {
	clearMenus()
	s := newServer()

	req := httptest.NewRequest(http.MethodDelete, "/menus/"+uuid.New().String(), nil)
	req.SetPathValue("id", uuid.New().String())
	w := httptest.NewRecorder()

	s.handleDeleteMenu(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleDeleteMenu_InvalidUUID_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodDelete, "/menus/bad", nil)
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	s.handleDeleteMenu(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleCreateDish ──────────────────────────────────────────────────────────

func TestHandleCreateDish_OK(t *testing.T) {
	clearRestaurants()
	clearMenus()
	clearDishes()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	s := newServer()

	payload := models.Dish{MenuID: m.Id, Name: "Pizza", Price: 12.5}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/dishes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleCreateDish(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Result().Body).Decode(&resp)
	if _, err := uuid.Parse(resp["id"]); err != nil {
		t.Errorf("expected valid UUID id, got: %s", resp["id"])
	}
}

func TestHandleCreateDish_MenuNotFound_Returns404(t *testing.T) {
	clearMenus()
	s := newServer()

	payload := models.Dish{MenuID: uuid.New(), Name: "Ghost dish"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/dishes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	s.handleCreateDish(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleCreateDish_MissingName_Returns400(t *testing.T) {
	clearRestaurants()
	clearMenus()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	s := newServer()

	payload := models.Dish{MenuID: m.Id}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/dishes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	s.handleCreateDish(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleCreateDish_InvalidJSON_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPost, "/dishes", bytes.NewReader([]byte(`{bad}`)))
	w := httptest.NewRecorder()

	s.handleCreateDish(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleUpdateDish ──────────────────────────────────────────────────────────

func TestHandleUpdateDish_OK(t *testing.T) {
	clearMenus()
	clearDishes()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	d := seedDish(m.Id)
	s := newServer()

	price := float32(19.99)
	payload := models.Dish{Name: "Updated Dish", Price: price}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/dishes/"+d.Id.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", d.Id.String())
	w := httptest.NewRecorder()

	s.handleUpdateDish(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	updated, _ := repository.DishSlice.GetByID(d.Id)
	if updated.Name != "Updated Dish" {
		t.Errorf("expected name 'Updated Dish', got '%s'", updated.Name)
	}
}

func TestHandleUpdateDish_NotFound_Returns404(t *testing.T) {
	clearDishes()
	s := newServer()

	body, _ := json.Marshal(models.Dish{Name: "x"})
	req := httptest.NewRequest(http.MethodPut, "/dishes/"+uuid.New().String(), bytes.NewReader(body))
	req.SetPathValue("id", uuid.New().String())
	w := httptest.NewRecorder()

	s.handleUpdateDish(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleUpdateDish_InvalidUUID_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPut, "/dishes/bad", bytes.NewReader([]byte(`{}`)))
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	s.handleUpdateDish(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleUpdateDish_InvalidJSON_Returns400(t *testing.T) {
	clearMenus()
	clearDishes()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	d := seedDish(m.Id)
	s := newServer()

	req := httptest.NewRequest(http.MethodPut, "/dishes/"+d.Id.String(), bytes.NewReader([]byte(`{bad}`)))
	req.SetPathValue("id", d.Id.String())
	w := httptest.NewRecorder()

	s.handleUpdateDish(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ── handleSaveCart ────────────────────────────────────────────────────────────

func TestHandleSaveCart_Open_OK(t *testing.T) {
	clearRestaurants()
	clearMenus()
	clearDishes()
	clearCarts()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	d := seedDish(m.Id)
	s := newServer()

	payload := models.Cart{
		RestaurantId: r.Id,
		DishIds:      []uuid.UUID{d.Id},
		Status:       models.CartStatusOpen,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleSaveCart(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Result().Body).Decode(&resp)
	if _, err := uuid.Parse(resp["id"]); err != nil {
		t.Errorf("expected valid UUID id, got: %s", resp["id"])
	}
}

func TestHandleSaveCart_Closed_CalculatesTotal(t *testing.T) {
	clearRestaurants()
	clearMenus()
	clearDishes()
	clearCarts()
	r := seedRestaurant()
	m := seedMenu(r.Id)
	d1 := seedDish(m.Id)
	d2 := seedDish(m.Id)
	s := newServer()

	payload := models.Cart{
		RestaurantId: r.Id,
		DishIds:      []uuid.UUID{d1.Id, d2.Id},
		Status:       models.CartStatusClosed,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleSaveCart(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]string
	json.NewDecoder(w.Result().Body).Decode(&resp)
	savedId, _ := uuid.Parse(resp["id"])
	cart, _ := repository.CartSlice.GetByID(savedId)
	if cart.Total != d1.Price+d2.Price {
		t.Errorf("expected total %.2f, got %.2f", d1.Price+d2.Price, cart.Total)
	}
}

func TestHandleSaveCart_RestaurantNotFound_Returns404(t *testing.T) {
	clearRestaurants()
	s := newServer()

	payload := models.Cart{RestaurantId: uuid.New(), Status: models.CartStatusOpen}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(body))
	w := httptest.NewRecorder()

	s.handleSaveCart(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleSaveCart_InvalidJSON_Returns400(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader([]byte(`{bad}`)))
	w := httptest.NewRecorder()

	s.handleSaveCart(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleSaveCart_DefaultStatusIsOpen(t *testing.T) {
	clearRestaurants()
	clearCarts()
	r := seedRestaurant()
	s := newServer()

	payload := models.Cart{RestaurantId: r.Id, DishIds: []uuid.UUID{}}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleSaveCart(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Result().Body).Decode(&resp)
	savedId, _ := uuid.Parse(resp["id"])
	cart, _ := repository.CartSlice.GetByID(savedId)
	if cart.Status != models.CartStatusOpen {
		t.Errorf("expected status 'open', got '%s'", cart.Status)
	}
}

// ── handleDeleteByIDRestaurants (cascade) ─────────────────────────────────────

func TestHandleDeleteRestaurant_CascadeDeletesMenusAndDishes(t *testing.T) {
	clearRestaurants()
	clearMenus()
	clearDishes()
	r := seedRestaurant()
	m1 := seedMenu(r.Id)
	m2 := seedMenu(r.Id)
	seedDish(m1.Id)
	seedDish(m1.Id)
	seedDish(m2.Id)
	s := newServer()

	req := httptest.NewRequest(http.MethodDelete, "/restaurants/"+r.Id.String(), nil)
	req.SetPathValue("id", r.Id.String())
	w := httptest.NewRecorder()

	s.handleDeleteByIDRestaurants(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if _, ok := repository.RestSlice.GetByID(r.Id); ok {
		t.Error("restaurant should have been deleted")
	}
	if repository.MenuSlice.Length() != 0 {
		t.Errorf("expected 0 menus after restaurant delete, got %d", repository.MenuSlice.Length())
	}
	if repository.DishSlice.Length() != 0 {
		t.Errorf("expected 0 dishes after restaurant delete, got %d", repository.DishSlice.Length())
	}
}
