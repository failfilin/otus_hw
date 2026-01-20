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
	Active   bool      `json:"active"`
}

func (res *Restaurant) ChangeActive() {
	switch res.Active {
	case true:
		res.Active = false
	case false:
		res.Active = true
	}
}

func (res *Restaurant) GetActive() bool {
	return res.Active
}

func (res Restaurant) String() string {
	//return fmt.Sprintf("ID заведения: %v \nНазвание: %v \nСсыль на лого: %v \nСписок менюшек: %v \nСтатус активности заведения: %t", res.Id, res.Name, res.Logo, res.MenuList, res.GetActive())
	return fmt.Sprintf("%v", res.Id) //заменил, глаза ломит на логах от него
}
func (r Restaurant) GetID() any {
	return r.Id
}

func (r *Restaurant) Update(name, logo *string, menuList *[]Menu, active *bool) {
	if name != nil {
		r.Name = *name
	}
	if logo != nil {
		r.Logo = *logo
	}
	if menuList != nil {
		r.MenuList = *menuList
	}
	if active != nil {
		r.ChangeActive()
	}
}
