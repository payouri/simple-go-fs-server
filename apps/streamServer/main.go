package main

import (
	"media-stream-server/config"
	"media-stream-server/routers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func createApp() {

	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*://*"},
	})
	app := echo.New()
	app.Use(cors)
	routers.InitRouter(app)

	serverAddress := ":" + config.EnvConfig.Port
	app.Logger.Fatal(app.Start(serverAddress))
}

func main() {
	createApp()
}
