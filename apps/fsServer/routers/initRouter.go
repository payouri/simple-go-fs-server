package routers

import (
	"fmt"
	"net/http"
	"simple-fs-web-service/auth"
	"simple-fs-web-service/routes"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var RouteMap = map[string]bool{}

func InitRouter(app *echo.Echo) {
	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*://*"},
	})
	app.Use(cors)
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			app.Logger.Print("Route called", c.Path())
			routeId := fmt.Sprintf("%s_%s", c.Request().Method, c.Path())
			startTime := time.Now()
			app.Logger.Printf("Route called: %s", routeId)
			if _, ok := RouteMap[routeId]; !ok {
				app.Logger.Printf("Route not found: %s", c.Path())
				c.String(http.StatusNotFound, "Not found")

				return nil
			}

			next(c)
			app.Logger.Print("Route executed", c.Path())
			app.Logger.Print("Time taken", startTime.Sub(startTime))

			return nil
		}
	}, auth.AuthenticationMiddleware)
	app.GET(routes.LIST_DIRECTORY_ROUTE_PATH, routes.ListEndpointRoute)
	app.GET(routes.DOWNLOAD_ROUTE_PATH, routes.DownloadEndpointRoute)
	app.GET(routes.GET_METADATA_ROUTE_PATH, routes.GetMetadataEndpointRoute)
	app.POST(routes.UPLOAD_ROUTE_PATH, routes.UploadEndpointRoute)

	app.RouteNotFound("*", func(c echo.Context) error {
		println("Path not found", c.Request().URL.Path)

		c.String(http.StatusNotFound, "Not found")

		return nil
	})

	for _, route := range app.Routes() {
		RouteMap[fmt.Sprintf("%s_%s", route.Method, route.Path)] = true
	}
}
