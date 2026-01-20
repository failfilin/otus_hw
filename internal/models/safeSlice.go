package models

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

type SafeSlice[T EatType] struct {
	Mu    sync.Mutex
	Slice []T
	File  string
}

func (safeSl *SafeSlice[T]) Append(eat T) {
	safeSl.Mu.Lock()
	defer safeSl.Mu.Unlock()
	safeSl.Slice = append(safeSl.Slice, eat)
}

func (safeSl *SafeSlice[T]) Length() int {
	safeSl.Mu.Lock()
	defer safeSl.Mu.Unlock()
	return len(safeSl.Slice)
}

func (safeSl *SafeSlice[T]) ReturnElement(i int) EatType {
	return safeSl.Slice[i]
}

func (safeSl *SafeSlice[T]) LogNewElements(state *int, printertype string) {
	if *state < safeSl.Length() {
		for i := *state; i < safeSl.Length(); i++ {
			log.Println("Новый тип", printertype, safeSl.Slice[i])
		}
		*state = safeSl.Length()
	}
}

func (s *SafeSlice[T]) GetByID(id uuid.UUID) (*T, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var zero T
	for i := range s.Slice {
		if s.Slice[i].GetID() == id {
			return &s.Slice[i], true // <-- возвращаем указатель на элемент внутри среза
		}
	}
	return &zero, false
}

func (s *SafeSlice[T]) RemoveByID(id uuid.UUID) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	for i, item := range s.Slice {
		if item.GetID() == id {
			s.Slice = append(s.Slice[:i], s.Slice[i+1:]...)
			return true
		}
	}
	return false
}
