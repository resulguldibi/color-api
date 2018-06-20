package service

import (
	"resulguldibi/color-api/repository"
)

type ColorService struct {
	colorRepository repository.ColorRepository
}

func NewColorService(colorRepository repository.ColorRepository) ColorService {
	return ColorService{colorRepository: colorRepository}
}
