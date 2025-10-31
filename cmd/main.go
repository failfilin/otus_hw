package main

import (
	"fmt"
	"os"
	"otus/repository"
	"otus/service"
)

func main() {

	var count int
	fmt.Println("Введи количество итераций")
	fmt.Fscan(os.Stdin, &count)
	for i := count; i > 0; i-- {
		repository.AddSlice(service.GenerateModels())
		fmt.Println("Вот и смотри куда добавило Рестораны", len(repository.RestSlice))
		fmt.Println("Вот и смотри куда добавило Меню", len(repository.MenuSlice))
		fmt.Println("Вот и смотри куда добавило Блюда", len(repository.DishSlice))
		fmt.Println("---------------------------------------")
	}
}
