package handler

import (
	"resulguldibi/color-api/service"
)

type ColorHandler struct {
	colorService service.ColorService
}

func NewColorHandler(colorService service.ColorService) ColorHandler {
	return ColorHandler{colorService: colorService}
}
