package repository

import (
	"log"
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
			RestSlice.LogNewElements(&restState, "ресторана")
			MenuSlice.LogNewElements(&menuState, "меню")
			DishSlice.LogNewElements(&dishState, "блюда")

			return

		case <-ticker.C:
			RestSlice.LogNewElements(&restState, "ресторан")
			MenuSlice.LogNewElements(&menuState, "меню")
			DishSlice.LogNewElements(&dishState, "блюдо")

		}

	}
}
