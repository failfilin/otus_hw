package models

type EatType interface {
	AddToSlice(slice *[]EatType)
}
