package routes

import (
	"log"
	"net/http"
	routesHelpers "simple-fs-web-service/routes/helpers"
	listDirectoriesRoute "simple-fs-web-service/routes/listDirectories"
	"simple-fs-web-service/validators"

	"github.com/labstack/echo/v4"
)

var LIST_DIRECTORY_ROUTE_PATH = routesHelpers.GetRoutePath("list/*")

func ListEndpointRoute(echoContext echo.Context) error {
	directoryPath := routesHelpers.GetFSFilePath(echoContext.Request().URL.Path, LIST_DIRECTORY_ROUTE_PATH)
	limitInt, limitErr := validators.ParseLimit(echoContext.QueryParam("limit"))
	offsetInt, offsetErr := validators.ParseOffset(echoContext.QueryParam("offset"))
	parseSortFieldResult, parseSortFieldErr := validators.ParseSortField(echoContext.QueryParam("sortField"))
	parseSortOrderResult, parseSortOrderErr := validators.ParseSortOrder(echoContext.QueryParam("sortOrder"))

	if limitErr != nil {
		log.Printf("Limit is not a valid integer: %s", limitErr.Error())
		echoContext.String(http.StatusBadRequest, "limit is not a valid integer")

		return nil
	}
	if offsetErr != nil {
		log.Printf("Offset is not a valid integer: %s", offsetErr.Error())
		echoContext.String(http.StatusBadRequest, "offset is not a valid integer")

		return nil
	}
	if parseSortFieldErr != nil {
		log.Printf("Sort field is not a valid: %s", parseSortFieldErr.Error())
		echoContext.String(http.StatusBadRequest, "sort field is not a valid")

		return nil
	}
	if parseSortOrderErr != nil {
		log.Printf("Sort order is not a valid: %s", parseSortOrderErr.Error())
		echoContext.String(http.StatusBadRequest, "sort order is not a valid")

		return nil
	}

	listDirectoryResult, listDirectoryError := listDirectoriesRoute.ListDirectory(listDirectoriesRoute.ListDirectoryParams{
		Path:      directoryPath,
		Limit:     limitInt,
		Offset:    offsetInt,
		SortField: parseSortFieldResult,
		SortOrder: parseSortOrderResult,
	})
	if listDirectoryError != nil {
		echoContext.String(http.StatusBadRequest, listDirectoryError.Error())
		log.Printf("Error while listing directory: %s", listDirectoryError.Error())

		return listDirectoryError
	}

	echoContext.JSON(200, map[string]interface{}{"files": listDirectoryResult.Files, "total": listDirectoryResult.Total, "maxPage": listDirectoryResult.MaxPage, "page": listDirectoryResult.Page})
	return nil
}
