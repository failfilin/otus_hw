package models

import (
	"log"
	"sync"
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
