package main

import (
	"fmt"
	"os"
	"otus/models"
	"otus/repository"
	"otus/service"
	"sync"
)

func main() {

	var count int
	var wgProduce sync.WaitGroup
	var wgConsume sync.WaitGroup
	var wgLogger sync.WaitGroup
	channel := make(chan models.EatType)
	doneChannel := make(chan struct{})

	fmt.Println("Введи количество итераций")
	fmt.Fscan(os.Stdin, &count)
	wgConsume.Add(1)
	go repository.AddSlice(channel, doneChannel, &wgConsume)
	wgLogger.Add(1)
	go repository.Logger(doneChannel, &wgLogger)
	for i := count; i > 0; i-- {
		wgProduce.Add(1)
		go service.GenerateModels(channel, &wgProduce)
	}
	wgProduce.Wait()
	close(channel)
	wgConsume.Wait()
	wgLogger.Wait()
}
