package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Restaurant struct {
	Id       uuid.UUID
	Name     string
	Logo     string // ссылка на CDN с картинкой
	MenuList []Menu
	active   bool
}

func (res *Restaurant) ChangeActive() {
	switch res.active {
	case true:
		res.active = false
	case false:
		res.active = true
	}
}

func (res *Restaurant) GetActive() bool {
	return res.active
}

func (res *Restaurant) ShowDetais() {
	fmt.Println("ID заведения", res.Id)
	fmt.Println("Название", res.Name)
	fmt.Println("Ссыль на лого", res.Logo)
	fmt.Println("Список менюшек", res.MenuList)
	fmt.Println("Статус активности заведения", res.GetActive())
}
