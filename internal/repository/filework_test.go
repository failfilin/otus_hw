package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/google/uuid"
)

// ── helpers ───────────────────────────────────────────────────────────────────

// tmpFile создаёт временный файл с переданным JSON-содержимым.
// Возвращает путь к файлу. Файл удаляется автоматически после теста.
func tmpFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	if content != "" {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write tmp file: %v", err)
		}
	}
	return path
}

// makeRestSlice создаёт SafeSlice[Restaurant] с указанным путём к файлу.
func makeRestSlice(file string) models.SafeSlice[models.Restaurant] {
	return models.SafeSlice[models.Restaurant]{File: file}
}

// sampleRestaurant возвращает тестовый ресторан с уникальным ID.
func sampleRestaurant() models.Restaurant {
	return models.Restaurant{
		Id:     uuid.New(),
		Name:   "Test Restaurant",
		Logo:   "logo.png",
		Active: true,
	}
}

// ── LoadFromFile ──────────────────────────────────────────────────────────────

// Сценарий: файл не существует — slice должен стать пустым, ошибки нет.
func TestLoadFromFile_FileNotExist_ReturnsEmptySlice(t *testing.T) {
	sl := makeRestSlice("/nonexistent/path/data.json")

	err := LoadFromFile(&sl)

	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if sl.Slice == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(sl.Slice) != 0 {
		t.Errorf("expected 0 elements, got %d", len(sl.Slice))
	}
}

// Сценарий: файл существует с валидным JSON — данные загружаются корректно.
func TestLoadFromFile_ValidJSON_LoadsData(t *testing.T) {
	rest := sampleRestaurant()
	data, _ := json.Marshal([]models.Restaurant{rest})
	path := tmpFile(t, string(data))

	sl := makeRestSlice(path)
	err := LoadFromFile(&sl)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sl.Slice) != 1 {
		t.Fatalf("expected 1 element, got %d", len(sl.Slice))
	}
	if sl.Slice[0].Id != rest.Id {
		t.Errorf("expected id %s, got %s", rest.Id, sl.Slice[0].Id)
	}
	if sl.Slice[0].Name != rest.Name {
		t.Errorf("expected name %q, got %q", rest.Name, sl.Slice[0].Name)
	}
}

// Сценарий: файл содержит несколько записей — все загружаются.
func TestLoadFromFile_MultipleRecords_LoadsAll(t *testing.T) {
	rests := []models.Restaurant{sampleRestaurant(), sampleRestaurant(), sampleRestaurant()}
	data, _ := json.Marshal(rests)
	path := tmpFile(t, string(data))

	sl := makeRestSlice(path)
	err := LoadFromFile(&sl)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sl.Slice) != 3 {
		t.Errorf("expected 3 elements, got %d", len(sl.Slice))
	}
}

// Сценарий: файл содержит невалидный JSON — возвращается ошибка декодирования.
func TestLoadFromFile_InvalidJSON_ReturnsError(t *testing.T) {
	path := tmpFile(t, `{not valid json}`)

	sl := makeRestSlice(path)
	err := LoadFromFile(&sl)

	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// Сценарий: файл содержит пустой JSON-массив — slice пустой, ошибки нет.
func TestLoadFromFile_EmptyJSONArray_ReturnsEmptySlice(t *testing.T) {
	path := tmpFile(t, `[]`)

	sl := makeRestSlice(path)
	err := LoadFromFile(&sl)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sl.Slice) != 0 {
		t.Errorf("expected 0 elements, got %d", len(sl.Slice))
	}
}

// Сценарий: повторный вызов LoadFromFile перезаписывает предыдущие данные в slice.
func TestLoadFromFile_CalledTwice_OverwritesPreviousData(t *testing.T) {
	rest1 := sampleRestaurant()
	rest2 := sampleRestaurant()

	data1, _ := json.Marshal([]models.Restaurant{rest1})
	path := tmpFile(t, string(data1))
	sl := makeRestSlice(path)
	_ = LoadFromFile(&sl)

	// Перезаписываем файл новыми данными
	data2, _ := json.Marshal([]models.Restaurant{rest2, rest2})
	_ = os.WriteFile(path, data2, 0644)

	err := LoadFromFile(&sl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sl.Slice) != 2 {
		t.Errorf("expected 2 elements after reload, got %d", len(sl.Slice))
	}
}

// ── SaveToFile ────────────────────────────────────────────────────────────────

// Сценарий: сохранение одного элемента — файл создаётся и содержит корректный JSON.
func TestSaveToFile_SingleElement_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	rest := sampleRestaurant()

	sl := makeRestSlice(path)
	sl.Slice = []models.Restaurant{rest}

	err := SaveToFile(&sl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected file to be created, but it does not exist")
	}
}

// Сценарий: сохранённый файл содержит те же данные, что были в slice.
func TestSaveToFile_DataRoundtrip_MatchesOriginal(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	rest := sampleRestaurant()

	sl := makeRestSlice(path)
	sl.Slice = []models.Restaurant{rest}
	_ = SaveToFile(&sl)

	// Читаем файл обратно и сравниваем
	sl2 := makeRestSlice(path)
	if err := LoadFromFile(&sl2); err != nil {
		t.Fatalf("failed to load saved file: %v", err)
	}
	if len(sl2.Slice) != 1 {
		t.Fatalf("expected 1 element, got %d", len(sl2.Slice))
	}
	if sl2.Slice[0].Id != rest.Id {
		t.Errorf("id mismatch: expected %s, got %s", rest.Id, sl2.Slice[0].Id)
	}
	if sl2.Slice[0].Name != rest.Name {
		t.Errorf("name mismatch: expected %q, got %q", rest.Name, sl2.Slice[0].Name)
	}
}

// Сценарий: сохранение пустого slice — файл создаётся с пустым JSON-массивом.
func TestSaveToFile_EmptySlice_WritesEmptyArray(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")

	sl := makeRestSlice(path)
	sl.Slice = []models.Restaurant{}
	err := SaveToFile(&sl)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, _ := os.ReadFile(path)
	var result []models.Restaurant
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("file content is not valid JSON: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty array in file, got %d elements", len(result))
	}
}

// Сценарий: вложенная директория не существует — SaveToFile создаёт её сам.
func TestSaveToFile_CreatesNestedDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "data.json")

	sl := makeRestSlice(path)
	sl.Slice = []models.Restaurant{sampleRestaurant()}

	err := SaveToFile(&sl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file in nested dir, but it does not exist")
	}
}

// Сценарий: повторный SaveToFile перезаписывает старые данные новыми.
func TestSaveToFile_Overwrite_ReplacesOldData(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")

	sl := makeRestSlice(path)
	sl.Slice = []models.Restaurant{sampleRestaurant(), sampleRestaurant()}
	_ = SaveToFile(&sl)

	// Перезаписываем одним элементом
	sl.Slice = []models.Restaurant{sampleRestaurant()}
	_ = SaveToFile(&sl)

	sl2 := makeRestSlice(path)
	_ = LoadFromFile(&sl2)
	if len(sl2.Slice) != 1 {
		t.Errorf("expected 1 element after overwrite, got %d", len(sl2.Slice))
	}
}

// ── InitRepository ────────────────────────────────────────────────────────────

// Сценарий: InitRepository вызывается когда файлов нет — глобальные слайсы
// инициализируются пустыми (не nil), без паники и без ошибок.
func TestInitRepository_FilesNotExist_SetsEmptySlices(t *testing.T) {
	// Подменяем пути на несуществующие, сохраняем оригиналы для восстановления
	origRest := RestSlice.File
	origMenu := MenuSlice.File
	origDish := DishSlice.File
	defer func() {
		RestSlice.File = origRest
		MenuSlice.File = origMenu
		DishSlice.File = origDish
	}()

	RestSlice.File = "/nonexistent/restaurants.json"
	MenuSlice.File = "/nonexistent/menus.json"
	DishSlice.File = "/nonexistent/dishes.json"

	// Не должно быть паники
	InitRepository()

	if RestSlice.Slice == nil {
		t.Error("RestSlice.Slice should not be nil after InitRepository")
	}
	if MenuSlice.Slice == nil {
		t.Error("MenuSlice.Slice should not be nil after InitRepository")
	}
	if DishSlice.Slice == nil {
		t.Error("DishSlice.Slice should not be nil after InitRepository")
	}
}

// Сценарий: InitRepository загружает данные из существующих файлов во все три слайса.
func TestInitRepository_WithFiles_LoadsAllSlices(t *testing.T) {
	dir := t.TempDir()

	// Готовим файлы для каждого слайса
	restPath := filepath.Join(dir, "restaurants.json")
	menuPath := filepath.Join(dir, "menus.json")
	dishPath := filepath.Join(dir, "dishes.json")

	rest := sampleRestaurant()
	menu := models.Menu{Id: uuid.New(), Description: "Test menu"}
	dish := models.Dish{Id: 1, Name: "Test dish"}

	writeJSON(t, restPath, []models.Restaurant{rest})
	writeJSON(t, menuPath, []models.Menu{menu})
	writeJSON(t, dishPath, []models.Dish{dish})

	origRest := RestSlice.File
	origMenu := MenuSlice.File
	origDish := DishSlice.File
	defer func() {
		RestSlice.File = origRest
		MenuSlice.File = origMenu
		DishSlice.File = origDish
	}()

	RestSlice.File = restPath
	MenuSlice.File = menuPath
	DishSlice.File = dishPath

	InitRepository()

	if len(RestSlice.Slice) != 1 {
		t.Errorf("expected 1 restaurant, got %d", len(RestSlice.Slice))
	}
	if len(MenuSlice.Slice) != 1 {
		t.Errorf("expected 1 menu, got %d", len(MenuSlice.Slice))
	}
	if len(DishSlice.Slice) != 1 {
		t.Errorf("expected 1 dish, got %d", len(DishSlice.Slice))
	}
}

// writeJSON — вспомогательная функция для записи JSON в файл в тестах.
func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}
