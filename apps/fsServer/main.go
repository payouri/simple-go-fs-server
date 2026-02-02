package main

import (
	"log"
	"simple-fs-web-service/auth"
	"simple-fs-web-service/config"
	"simple-fs-web-service/routers"

	"github.com/labstack/echo/v4"
)

func createApp() {

	if !auth.ApiKeyStore.HasOneExistingKey() {
		_, generateKeyError := auth.ApiKeyStore.GenerateApiKey()

		if generateKeyError != nil {
			log.Fatal(generateKeyError)
		}
	}

	app := echo.New()
	routers.InitRouter(app)

	serverAddress := ":" + config.EnvConfig.Port
	app.Logger.Fatal(app.Start(serverAddress))
}

func main() {
	createApp()
}
