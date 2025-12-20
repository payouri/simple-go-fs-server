package main

import (
	"log"
	"os"
	"simple-fs-web-service/auth"
	"simple-fs-web-service/routers"

	"github.com/labstack/echo/v4"
)

func createApp() {
	envPort := os.Getenv("FS_SERVER_PORT")
	if envPort == "" {
		envPort = "5008"
	}
	if !auth.ApiKeyStore.HasOneExistingKey() {
		_, generateKeyError := auth.ApiKeyStore.GenerateApiKey()

		if generateKeyError != nil {
			log.Fatal(generateKeyError)
		}
	}

	app := echo.New()
	routers.InitRouter(app)

	serverAddress := ":" + envPort
	app.Logger.Fatal(app.Start(serverAddress))
}

func main() {
	createApp()
}
