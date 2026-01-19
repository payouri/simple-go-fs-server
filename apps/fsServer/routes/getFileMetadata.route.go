package routes

import (
	"net/http"
	getFileMetadataRoute "simple-fs-web-service/routes/getFileMetadata"
	routesHelpers "simple-fs-web-service/routes/helpers"

	"github.com/labstack/echo/v4"
)

var GET_METADATA_ROUTE_PATH = routesHelpers.GetRoutePath("metadata/*")

func GetMetadataEndpointRoute(echoContext echo.Context) error {
	filePath := routesHelpers.GetFSFilePath(echoContext.Request().URL.Path, GET_METADATA_ROUTE_PATH)

	metadataResult, metadataError := getFileMetadataRoute.GetFileMetadata(getFileMetadataRoute.GetFileMetadataParams{
		FilePath: filePath,
	})
	if metadataError != nil {
		echoContext.String(http.StatusBadRequest, metadataError.Error())

		return metadataError
	}

	return echoContext.JSON(http.StatusOK, metadataResult)
}
