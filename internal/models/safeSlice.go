package models

import (
	"log"
	"sync"
)

type SafeSlice[slice EatType] struct {
	mu    sync.Mutex
	slice []slice
}

func (safeSl *SafeSlice[EatType]) Append(eat EatType) {
	safeSl.mu.Lock()
	defer safeSl.mu.Unlock()
	safeSl.slice = append(safeSl.slice, eat)
}

func (safeSl *SafeSlice[EatType]) Length() int {
	safeSl.mu.Lock()
	defer safeSl.mu.Unlock()
	return len(safeSl.slice)
}

func (safeSl *SafeSlice[EatType]) ReturnElement(i int) EatType {
	return safeSl.slice[i]
}

func (safeSl *SafeSlice[EatType]) LogNewElements(state *int, printertype string) {
	if *state < safeSl.Length() {
		for i := *state; i < safeSl.Length(); i++ {
			log.Println("Новый тип", printertype, safeSl.slice[i])
		}
		*state = safeSl.Length()
	}
}
