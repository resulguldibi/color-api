package main

import (
	"resulguldibi/color-api/server"
)

func main() {
	server.NewServer().Run(":8080")
}
