package routes

import "github.com/labstack/echo/v4"

func StreamAudioFileRoute(echoContext echo.Context) error {
	echoContext.String(200, "Hello World")

	return nil
}
