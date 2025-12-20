package main

import (
	"media-stream-server/config"
	"media-stream-server/routers"

	"github.com/labstack/echo/v4"
)

func createApp() {

	app := echo.New()
	routers.InitRouter(app)

	serverAddress := ":" + config.EnvConfig.Port
	app.Logger.Fatal(app.Start(serverAddress))
}

func main() {
	createApp()
}
