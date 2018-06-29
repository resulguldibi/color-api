package server

import (
	"resulguldibi/color-api/factory"
	"resulguldibi/color-api/handler"
	"resulguldibi/color-api/middleware"
	"resulguldibi/color-api/repository"
	"resulguldibi/color-api/service"

	httpClientFactory "resulguldibi/http-client/factory"
	redisClientFactory "resulguldibi/redis-client/factory"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	factory.InitFactoryList()
	AddDefaultMiddlewaresToEngine(server)
	//TODO : get connection info from config
	dbClientFactory := repository.NewDbClientFactory("sqlite3", "./SQLiteDB.db")
	redisClientFactory := redisClientFactory.NewRedisClientFactory("localhost:6379", "")

	httpClientFactory := httpClientFactory.NewHttpClientFactory()

	server.GET("/colors", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorService(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient()))
		colorHandler.HandleGetRandomColors(ctx)
	})

	server.GET("/levels", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorService(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient()))
		colorHandler.HandleGetLevels(ctx)
	})

	server.POST("/validate", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorService(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient()))
		colorHandler.HandleValidateColors(ctx)
	})

	server.POST("/google/oauth2/token", func(ctx *gin.Context) {
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		userHandler.HandleOAuth2Google(ctx)
	})

	server.GET("/ranking", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorService(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient()))
		colorHandler.HandleRankings(ctx)
	})

	return server
}

func AddDefaultMiddlewaresToEngine(server *gin.Engine) {
	//engine.Use(secure.Secure(secure.Options))
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.Use(middleware.UseUserMiddleware())
}
