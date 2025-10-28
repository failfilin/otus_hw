package models

type EatType interface {
	GenerateNewOne() EatType
}
