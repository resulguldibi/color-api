package handler

import (
	"resulguldibi/color-api/service"
)

type ColorHandler struct {
	colorService service.ColorService
}

type UserHandler struct {
	userService service.UserService
}

type SocketHandler struct{
	socketService service.SocketService
}

func NewColorHandler(colorService service.ColorService) ColorHandler {
	return ColorHandler{colorService: colorService}
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService: userService}
}

func NewSocketHandler(socketService service.SocketService) SocketHandler {
	return SocketHandler{socketService: socketService}
}
