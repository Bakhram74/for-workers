package main

import (
	"github.com/ShamilKhal/shgo/config"
	_ "github.com/ShamilKhal/shgo/docs"
	"github.com/ShamilKhal/shgo/internal/app"
	"github.com/ShamilKhal/shgo/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Shupir app API
// @version         1.0
// @description     API docs for Shupir Application.

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	logger.InitLog(cfg.LogLevel)

	go swaggerRun()

	app.Run(cfg)
}

func swaggerRun() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Log.Info().Msgf("Run Swagger at address %s", ":8081/swagger/index.html#/")
	r.Run(":8081")
}
