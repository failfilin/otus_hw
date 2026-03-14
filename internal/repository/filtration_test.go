package repository

import (
	"testing"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/google/uuid"
)

// clearAll обнуляет все три глобальных слайса перед тестом,
// чтобы данные из предыдущих тестов не влияли на текущий.
func clearAll() {
	RestSlice.Mu.Lock()
	RestSlice.Slice = nil
	RestSlice.Mu.Unlock()

	MenuSlice.Mu.Lock()
	MenuSlice.Slice = nil
	MenuSlice.Mu.Unlock()

	DishSlice.Mu.Lock()
	DishSlice.Slice = nil
	DishSlice.Mu.Unlock()
}

// Сценарий: передаём Restaurant — он попадает в RestSlice.
func TestAddSlice_Restaurant_AppendsToRestSlice(t *testing.T) {
	clearAll()

	rest := models.Restaurant{
		Id:   uuid.New(),
		Name: "New Restaurant",
	}

	AddSlice(rest)

	if RestSlice.Length() != 1 {
		t.Errorf("expected 1 restaurant in RestSlice, got %d", RestSlice.Length())
	}
}

// Сценарий: передаём Menu — он попадает в MenuSlice.
func TestAddSlice_Menu_AppendsToMenuSlice(t *testing.T) {
	clearAll()

	menu := models.Menu{
		Id:          uuid.New(),
		Description: "Lunch menu",
	}

	AddSlice(menu)

	if MenuSlice.Length() != 1 {
		t.Errorf("expected 1 menu in MenuSlice, got %d", MenuSlice.Length())
	}
}

// Сценарий: передаём Dish — он попадает в DishSlice.
func TestAddSlice_Dish_AppendsToDishSlice(t *testing.T) {
	clearAll()

	dish := models.Dish{
		Id:   uuid.New(),
		Name: "Pasta",
	}

	AddSlice(dish)

	if DishSlice.Length() != 1 {
		t.Errorf("expected 1 dish in DishSlice, got %d", DishSlice.Length())
	}
}

// Сценарий: AddSlice не затрагивает слайсы других типов.
func TestAddSlice_Restaurant_DoesNotAffectOtherSlices(t *testing.T) {
	clearAll()

	AddSlice(models.Restaurant{Id: uuid.New(), Name: "Solo"})

	if MenuSlice.Length() != 0 {
		t.Errorf("MenuSlice should be untouched, got %d elements", MenuSlice.Length())
	}
	if DishSlice.Length() != 0 {
		t.Errorf("DishSlice should be untouched, got %d elements", DishSlice.Length())
	}
}

// Сценарий: несколько вызовов AddSlice накапливают элементы в слайсе.
func TestAddSlice_MultipleRestaurants_AccumulatesAll(t *testing.T) {
	clearAll()

	AddSlice(models.Restaurant{Id: uuid.New(), Name: "First"})
	AddSlice(models.Restaurant{Id: uuid.New(), Name: "Second"})
	AddSlice(models.Restaurant{Id: uuid.New(), Name: "Third"})

	if RestSlice.Length() != 3 {
		t.Errorf("expected 3 restaurants, got %d", RestSlice.Length())
	}
}
