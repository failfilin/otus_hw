package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/failfilin/otus_hw/internal/repository"
)

func Logger(ctx context.Context, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	restState := 0
	dishState := 0
	menuState := 0

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
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
