package main

import (
	"fmt"
	"os"
	"otus/models"
	"otus/repository"
	"otus/service"
)

func main() {
	var RestSlice []models.EatType
	var MenuSlice []models.EatType
	var DishSlice []models.EatType
	var count int
	fmt.Println("Введи количество итераций")
	fmt.Fscan(os.Stdin, &count)
	for i := count; i > 0; i-- {
		repository.AddSlice(service.GenerateModels(), &DishSlice, &MenuSlice, &RestSlice)
		fmt.Println("Вот и смотри куда добавило Рестораны", len(RestSlice))
		fmt.Println("Вот и смотри куда добавило Меню", len(MenuSlice))
		fmt.Println("Вот и смотри куда добавило Блюда", len(DishSlice))
		fmt.Println("---------------------------------------")
	}
}
