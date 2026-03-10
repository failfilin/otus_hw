package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/failfilin/otus_hw/internal/repository"
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

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
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
