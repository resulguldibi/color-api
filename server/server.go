package server

import (
	"net/http"
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

	server.LoadHTMLGlob("static/html/*.html")
	server.Static("/css", "static/css")
	server.Static("/js", "static/js")
	server.Static("/images", "static/images")
	//server.Use(static.Serve("/assets", static.LocalFile("/assets", false)))

	hub := service.NewSocketHub(redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient())
	go hub.Broadcast()
	go hub.Register()
	go hub.UnRegister()
	go hub.RegisterMatch()
	go hub.UnRegisterMatch()
	go hub.RegisterMultiPlay()
	go hub.AcceptMatchForMultiPlay()
	go hub.UnRegisterMultiPlay()
	go hub.MultiPlayMatchMove()
	go hub.MultiPlayMatchMessage()

	server.GET("/google/oauth2", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "googleoauth2.html", nil)
	})

	server.GET("/play", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "play.html", nil)
	})

	server.GET("/multiplay", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "multiplay.html", nil)
	})

	server.GET("/stage", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "stage.html", nil)
	})

	server.GET("/help", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		colorHandler.HandleColorHelp(ctx)
	})

	server.GET("/stephelp", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		colorHandler.HandleColorStepHelp(ctx)
	})

	server.GET("/colors", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		colorHandler.HandleGetRandomColors(ctx)
	})

	server.GET("/history/raund", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		colorHandler.HandleRaundHistory(ctx)
	})

	server.GET("/user/stage", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		colorHandler.HandleGetUserStageInfo(ctx)
	})

	server.GET("/levels", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		colorHandler.HandleGetLevels(ctx)
	})

	server.POST("/validate", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		colorHandler.HandleValidateColors(ctx)
	})

	server.POST("/google/oauth2/token", func(ctx *gin.Context) {
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		userHandler.HandleOAuth2Google(ctx)
	})

	server.GET("/multiplay/register", func(ctx *gin.Context) {
		socketHandler := handler.NewSocketHandler(service.NewSocketService())
		socketHandler.HandleRegisterForMultiPlay(ctx, hub)
	})

	server.GET("/multiplay/accept", func(ctx *gin.Context) {
		socketHandler := handler.NewSocketHandler(service.NewSocketService())
		socketHandler.HandleAcceptMatchForMultiPlay(ctx, hub)
	})

	server.GET("/multiplay/unregister", func(ctx *gin.Context) {
		socketHandler := handler.NewSocketHandler(service.NewSocketService())
		socketHandler.HandleUnRegisterForMultiPlay(ctx, hub)
	})

	server.POST("/multiplay/move", func(ctx *gin.Context) {
		socketHandler := handler.NewSocketHandler(service.NewSocketService())
		socketHandler.HandleMultiplayMove(ctx, hub)
	})

	server.POST("/multiplay/validate", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
		colorHandler.HandleMultiPlayValidateColors(ctx, hub)
	})

	server.GET("/ranking", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		colorHandler := handler.NewColorHandler(service.NewColorServiceHttpClient(repository.NewColorRepository(dbClient), redisClientFactory.GetRedisClient(), httpClientFactory.GetHttpClient()))
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
