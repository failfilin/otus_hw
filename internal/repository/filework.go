package repository

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/failfilin/otus_hw/internal/models"
)

func LoadFromFile[T models.EatType](safeSl *models.SafeSlice[T]) error {
	safeSl.Mu.Lock()
	defer safeSl.Mu.Unlock()

	file, err := os.Open(safeSl.File)
	if err != nil {
		if os.IsNotExist(err) { // файла нет — просто оставляем slice пустым
			safeSl.Slice = []T{}
			return nil
		}
		return err
	}
	defer file.Close()

	var data []T
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	safeSl.Slice = data
	return nil
}

func SaveToFile[T models.EatType](safeSl *models.SafeSlice[T]) error {
	if err := os.MkdirAll(filepath.Dir(safeSl.File), 0755); err != nil {
		return err
	}
	file, err := os.Create(safeSl.File)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(safeSl.Slice)
}

func InitRepository() {
	_ = LoadFromFile(&RestSlice)
	_ = LoadFromFile(&MenuSlice)
	_ = LoadFromFile(&DishSlice)
}
