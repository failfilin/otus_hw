package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Restaurant struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Logo     string    `json:"logo"` // ссылка на CDN с картинкой
	MenuList []Menu    `json:"menuList"`
	active   bool      `json:"active"`
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

func (res Restaurant) String() string {
	//return fmt.Sprintf("ID заведения: %v \nНазвание: %v \nСсыль на лого: %v \nСписок менюшек: %v \nСтатус активности заведения: %t", res.Id, res.Name, res.Logo, res.MenuList, res.GetActive())
	return fmt.Sprintf("%v", res.Id) //заменил, глаза ломит на логах от него
}
