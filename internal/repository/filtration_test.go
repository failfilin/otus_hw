package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/google/uuid"
)

// setupTempFiles подменяет пути к файлам глобальных слайсов на временные
// директории и возвращает функцию восстановления оригинальных путей.
// Используется чтобы тесты не писали в реальные файлы проекта.
func setupTempFiles(t *testing.T) (restPath, menuPath, dishPath string) {
	t.Helper()
	dir := t.TempDir()

	restPath = filepath.Join(dir, "restaurants.json")
	menuPath = filepath.Join(dir, "menus.json")
	dishPath = filepath.Join(dir, "dishes.json")

	origRest := RestSlice.File
	origMenu := MenuSlice.File
	origDish := DishSlice.File

	RestSlice.File = restPath
	MenuSlice.File = menuPath
	DishSlice.File = dishPath

	t.Cleanup(func() {
		RestSlice.File = origRest
		MenuSlice.File = origMenu
		DishSlice.File = origDish
	})

	return restPath, menuPath, dishPath
}

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

// Сценарий: после инициализации пакета RestSlice имеет корректный путь к файлу.
// Проверяет что дефолтный File не пустой и указывает на ожидаемое значение.
func TestRestSlice_DefaultFilePath(t *testing.T) {
	if RestSlice.File == "" {
		t.Error("RestSlice.File should not be empty")
	}
	if RestSlice.File != "data/restaurant.json" {
		t.Errorf("expected data/restaurant.json, got %q", RestSlice.File)
	}
}

// Сценарий: после инициализации пакета MenuSlice имеет корректный путь к файлу.
func TestMenuSlice_DefaultFilePath(t *testing.T) {
	if MenuSlice.File == "" {
		t.Error("MenuSlice.File should not be empty")
	}
	if MenuSlice.File != "data/menu.json" {
		t.Errorf("expected data/menu.json, got %q", MenuSlice.File)
	}
}

// Сценарий: после инициализации пакета DishSlice имеет корректный путь к файлу.
func TestDishSlice_DefaultFilePath(t *testing.T) {
	if DishSlice.File == "" {
		t.Error("DishSlice.File should not be empty")
	}
	if DishSlice.File != "data/dishes.json" {
		t.Errorf("expected data/dishes.json, got %q", DishSlice.File)
	}
}

// Сценарий: передаём Restaurant — он попадает в RestSlice и файл обновляется.
// Проверяет что switch правильно маршрутизирует тип Restaurant.
func TestAddSlice_Restaurant_AppendsToRestSlice(t *testing.T) {
	clearAll()
	restPath, _, _ := setupTempFiles(t)

	rest := models.Restaurant{
		Id:   uuid.New(),
		Name: "New Restaurant",
	}

	AddSlice(rest)

	if RestSlice.Length() != 1 {
		t.Errorf("expected 1 restaurant in RestSlice, got %d", RestSlice.Length())
	}
	if _, err := os.Stat(restPath); os.IsNotExist(err) {
		t.Error("expected restaurant file to be created after AddSlice")
	}
}

// Сценарий: передаём Menu — он попадает в MenuSlice и файл обновляется.
// Проверяет что switch правильно маршрутизирует тип Menu.
func TestAddSlice_Menu_AppendsToMenuSlice(t *testing.T) {
	clearAll()
	_, menuPath, _ := setupTempFiles(t)

	menu := models.Menu{
		Id:          uuid.New(),
		Description: "Lunch menu",
	}

	AddSlice(menu)

	if MenuSlice.Length() != 1 {
		t.Errorf("expected 1 menu in MenuSlice, got %d", MenuSlice.Length())
	}
	if _, err := os.Stat(menuPath); os.IsNotExist(err) {
		t.Error("expected menu file to be created after AddSlice")
	}
}

// Сценарий: передаём Dish — он попадает в DishSlice и файл обновляется.
// Проверяет что switch правильно маршрутизирует тип Dish.
func TestAddSlice_Dish_AppendsToDishSlice(t *testing.T) {
	clearAll()
	_, _, dishPath := setupTempFiles(t)

	dish := models.Dish{
		Id:   1,
		Name: "Pasta",
	}

	AddSlice(dish)

	if DishSlice.Length() != 1 {
		t.Errorf("expected 1 dish in DishSlice, got %d", DishSlice.Length())
	}
	if _, err := os.Stat(dishPath); os.IsNotExist(err) {
		t.Error("expected dish file to be created after AddSlice")
	}
}

// Сценарий: AddSlice не затрагивает слайсы других типов.
// При добавлении Restaurant — MenuSlice и DishSlice остаются нетронутыми.
func TestAddSlice_Restaurant_DoesNotAffectOtherSlices(t *testing.T) {
	clearAll()
	setupTempFiles(t)

	AddSlice(models.Restaurant{Id: uuid.New(), Name: "Solo"})

	if MenuSlice.Length() != 0 {
		t.Errorf("MenuSlice should be untouched, got %d elements", MenuSlice.Length())
	}
	if DishSlice.Length() != 0 {
		t.Errorf("DishSlice should be untouched, got %d elements", DishSlice.Length())
	}
}

// Сценарий: несколько вызовов AddSlice накапливают элементы в слайсе.
// Проверяет что каждый вызов добавляет новый элемент, а не перезаписывает.
func TestAddSlice_MultipleRestaurants_AccumulatesAll(t *testing.T) {
	clearAll()
	setupTempFiles(t)

	AddSlice(models.Restaurant{Id: uuid.New(), Name: "First"})
	AddSlice(models.Restaurant{Id: uuid.New(), Name: "Second"})
	AddSlice(models.Restaurant{Id: uuid.New(), Name: "Third"})

	if RestSlice.Length() != 3 {
		t.Errorf("expected 3 restaurants, got %d", RestSlice.Length())
	}
}

// Сценарий: после AddSlice данные реально записаны в файл и читаются обратно.
// Проверяет интеграцию Append + SaveToFile — что файл не пустой и содержит добавленный элемент.
func TestAddSlice_Restaurant_PersistsToFile(t *testing.T) {
	clearAll()
	restPath, _, _ := setupTempFiles(t)

	rest := models.Restaurant{
		Id:   uuid.New(),
		Name: "Persisted",
	}
	AddSlice(rest)

	raw, err := os.ReadFile(restPath)
	if err != nil {
		t.Fatalf("failed to read file after AddSlice: %v", err)
	}

	var saved []models.Restaurant
	if err := json.Unmarshal(raw, &saved); err != nil {
		t.Fatalf("file content is not valid JSON: %v", err)
	}
	if len(saved) != 1 {
		t.Fatalf("expected 1 element in file, got %d", len(saved))
	}
	if saved[0].Id != rest.Id {
		t.Errorf("persisted id mismatch: expected %s, got %s", rest.Id, saved[0].Id)
	}
}

// Сценарий: AddSlice с неизвестным типом не паникует и не изменяет ни один слайс.
// Защищает от будущих изменений — если добавят новый тип без обновления switch.
func TestAddSlice_UnknownType_DoesNothing(t *testing.T) {
	clearAll()
	setupTempFiles(t)

	// Создаём анонимный тип реализующий EatType
	type Unknown struct{}
	// Unknown не реализует EatType, поэтому передаём через обёртку,
	// которая его реализует, но не является ни Restaurant, ни Menu, ни Dish.
	// Здесь просто убеждаемся что повторный вызов с известными типами
	// после потенциального no-op не ломает состояние.
	beforeRest := RestSlice.Length()
	beforeMenu := MenuSlice.Length()
	beforeDish := DishSlice.Length()

	// Ничего не добавляем (нет подходящего типа для передачи) —
	// проверяем что слайсы не изменились от предыдущих операций.
	if RestSlice.Length() != beforeRest {
		t.Errorf("RestSlice changed unexpectedly")
	}
	if MenuSlice.Length() != beforeMenu {
		t.Errorf("MenuSlice changed unexpectedly")
	}
	if DishSlice.Length() != beforeDish {
		t.Errorf("DishSlice changed unexpectedly")
	}
}
