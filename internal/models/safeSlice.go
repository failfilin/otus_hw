package models

import (
	"sync"

	"github.com/google/uuid"
)

type SafeSlice[T EatType] struct {
	Mu    sync.Mutex
	Slice []T
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

func (s *SafeSlice[T]) GetByID(id uuid.UUID) (*T, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var zero T
	for i := range s.Slice {
		if s.Slice[i].GetID() == id {
			return &s.Slice[i], true
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

func (s *SafeSlice[T]) RemoveWhere(predicate func(T) bool) int {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	removed := 0
	result := make([]T, 0, len(s.Slice))
	for _, item := range s.Slice {
		if predicate(item) {
			removed++
		} else {
			result = append(result, item)
		}
	}
	s.Slice = result
	return removed
}
