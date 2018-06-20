package server

import (
	"resulguldibi/color-api/factory"
	"resulguldibi/color-api/handler"
	"resulguldibi/color-api/repository"
	"resulguldibi/color-api/service"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	factory.InitFactoryList()
	AddDefaultMiddlewaresToEngine(server)
	//TODO : get connection info from config
	dbClientFactory := repository.NewDbClientFactory("sqlite3", "./SQLiteDB.db")

	server.GET("/colors", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorService(repository.NewColorRepository(dbClient)))
		colorHandler.HandleGetRandomColors(ctx)
	})

	server.POST("/validate", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorService(repository.NewColorRepository(dbClient)))
		colorHandler.HandleValidateColors(ctx)
	})

	return server
}

func AddDefaultMiddlewaresToEngine(server *gin.Engine) {
	//engine.Use(secure.Secure(secure.Options))
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
}
