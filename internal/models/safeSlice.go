package models

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type SafeSlice[T EatType] struct {
	mu    sync.Mutex
	slice []T
	File  string
}

func (safeSl *SafeSlice[T]) Append(eat T) error {
	safeSl.mu.Lock()
	defer safeSl.mu.Unlock()
	safeSl.slice = append(safeSl.slice, eat)
	return safeSl.saveToFile()
}

func (safeSl *SafeSlice[T]) saveToFile() error {
	file, err := os.Create(safeSl.File)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(safeSl.slice)
}

func (safeSl *SafeSlice[T]) Length() int {
	safeSl.mu.Lock()
	defer safeSl.mu.Unlock()
	return len(safeSl.slice)
}

func (safeSl *SafeSlice[T]) ReturnElement(i int) EatType {
	return safeSl.slice[i]
}

func (safeSl *SafeSlice[T]) LogNewElements(state *int, printertype string) {
	if *state < safeSl.Length() {
		for i := *state; i < safeSl.Length(); i++ {
			log.Println("Новый тип", printertype, safeSl.slice[i])
		}
		*state = safeSl.Length()
	}
}

func (safeSl *SafeSlice[T]) LoadFromFile() error {
	safeSl.mu.Lock()
	defer safeSl.mu.Unlock()

	file, err := os.Open(safeSl.File)
	if err != nil {
		if os.IsNotExist(err) { // файла нет — просто оставляем slice пустым
			safeSl.slice = []T{}
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

	safeSl.slice = data
	return nil
}
