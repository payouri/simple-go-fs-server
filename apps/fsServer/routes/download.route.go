package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"simple-fs-web-service/helpers"
	routesHelpers "simple-fs-web-service/routes/helpers"

	"github.com/labstack/echo/v4"
)

var DOWNLOAD_ROUTE_PATH = routesHelpers.GetRoutePath("download/*")

func DownloadEndpointRoute(echoContext echo.Context) error {
	resolvedFsPath, resolvePathError := helpers.ResolvePath(helpers.GetBasePath(), routesHelpers.GetFSFilePath(echoContext.Request().URL.Path, DOWNLOAD_ROUTE_PATH))
	if resolvePathError != nil {
		log.Printf("Error while resolving path: %s", resolvePathError.Error())
		echoContext.String(http.StatusBadRequest, "Unable to resolve path")

		return resolvePathError
	}
	isPathLegalResult, isPathLegalError := helpers.IsPathLegal(helpers.GetBasePath(), resolvedFsPath)
	if isPathLegalError != nil {
		log.Printf("Error while checking path: %s", isPathLegalError.Error())
		echoContext.String(http.StatusBadRequest, "Unable to check path")

		return isPathLegalError
	}
	if !isPathLegalResult {
		log.Printf("Path is not a valid path")
		echoContext.String(http.StatusBadRequest, "Path is not a valid path")

		return nil
	}
	file, fileError := os.Open(resolvedFsPath)
	fileStat, fileStatError := file.Stat()

	if fileError != nil {
		log.Printf("Error while opening file: %s", fileError.Error())
		echoContext.String(http.StatusBadRequest, fileError.Error())

		return fileError
	}
	if fileStatError != nil {
		log.Printf("Error while getting file stat: %s", fileStatError.Error())
		echoContext.String(http.StatusBadRequest, fileStatError.Error())

		return fileStatError
	}

	defer file.Close()

	echoContext.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path.Base(resolvedFsPath)))
	echoContext.Response().Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))
	return echoContext.Stream(http.StatusOK, "application/octet-stream", file)
}
