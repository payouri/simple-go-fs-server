package routers

import (
	"fmt"
	"media-stream-server/routes"
	"net/http"

	"github.com/labstack/echo/v4"
)

var RouteMap = map[string]bool{}

func InitRouter(app *echo.Echo) {
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			app.Logger.Print("Route called", c.Path())
			routeId := fmt.Sprintf("%s_%s", c.Request().Method, c.Path())
			if _, ok := RouteMap[routeId]; !ok {
				app.Logger.Print("Route not found", c.Path())
				c.String(http.StatusNotFound, "Not found")

				return nil
			}

			next(c)

			return nil
		}
	})

	app.GET(routes.STREAM_AUDIO_FILE_ROUTE_PATH, routes.StreamAudioFileRoute)

	app.RouteNotFound("*", func(c echo.Context) error {
		println("Path not found", c.Request().URL.Path)

		c.String(http.StatusNotFound, "Not found")

		return nil
	})

	for _, route := range app.Routes() {
		RouteMap[fmt.Sprintf("%s_%s", route.Method, route.Path)] = true
	}
}
