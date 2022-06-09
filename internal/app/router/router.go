package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/basic/docs"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/app/handler"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
)

var Router *gin.Engine

func SetupRouter() (router *gin.Engine) {
	docs.SwaggerInfo.Version = config.Env.Version

	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router = gin.Default()
	router.Use(handler.CORSMiddleware(), handler.ErrorMiddleware())

	router.GET("/health", handler.HealthHandler)
	router.GET("/version", handler.VersionHandler)

	return
}

func Setup() error {
	Router = SetupRouter()

	return nil
}
