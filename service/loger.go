package service

import (
	"log"
	"otus/repository"
	"sync"
	"time"
)

func Logger(done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	restState := 0
	dishState := 0
	menuState := 0

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			// Финальный лог по завершении
			log.Println("Финальное срабатывание")
			repository.RestSlice.LogNewElements(&restState, "ресторана")
			repository.MenuSlice.LogNewElements(&menuState, "меню")
			repository.DishSlice.LogNewElements(&dishState, "блюда")

			return

		case <-ticker.C:
			repository.RestSlice.LogNewElements(&restState, "ресторан")
			repository.MenuSlice.LogNewElements(&menuState, "меню")
			repository.DishSlice.LogNewElements(&dishState, "блюдо")

		}

	}
}
