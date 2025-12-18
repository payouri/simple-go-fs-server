package routes

import (
	"log"
	"net/http"
	"simple-fs-web-service/helpers"
	routesHelpers "simple-fs-web-service/routes/helpers"
	uploadfile "simple-fs-web-service/routes/uploadFile"

	"github.com/labstack/echo/v4"
)

var UPLOAD_ROUTE_PATH = routesHelpers.GetRoutePath("upload/*")

func UploadEndpointRoute(echoContext echo.Context) error {
	log.Printf("Upload endpoint called")
	resolvedFsPath, resolvePathError := helpers.ResolvePath(helpers.GetBasePath(), routesHelpers.GetFSFilePath(echoContext.Request().URL.Path, UPLOAD_ROUTE_PATH))
	if resolvePathError != nil {
		log.Printf("Error while resolving path: %s", resolvePathError.Error())
		echoContext.String(http.StatusBadRequest, "Unable to resolve upload path")

		return resolvePathError
	}
	isPathLegalResult, isPathLegalError := helpers.IsPathLegal(helpers.GetBasePath(), resolvedFsPath)
	if isPathLegalError != nil {
		log.Printf("Error while checking path: %s", isPathLegalError.Error())
		echoContext.String(http.StatusBadRequest, "Unable to check upload path")

		return isPathLegalError
	}
	if !isPathLegalResult {
		log.Printf("Upload path is not a valid path")
		echoContext.String(http.StatusBadRequest, "Upload path is not a valid path")

		return nil
	}

	fileHeader, fileError := echoContext.FormFile("file")
	if fileError != nil {
		log.Printf("Error while getting file: %s", fileError.Error())
		echoContext.String(http.StatusBadRequest, "Unable to process file")

		return fileError
	}

	if fileHeader.Size > routesHelpers.MAX_UPLOAD_SIZE {
		log.Printf("File too big: %vMB", fileHeader.Size/1024/1024)
		echoContext.String(http.StatusRequestEntityTooLarge, "File too big")

		return nil
	}

	file, openFileError := fileHeader.Open()
	if openFileError != nil {
		log.Printf("Error while opening file: %s", openFileError.Error())
		echoContext.String(http.StatusBadRequest, "Unable to process file")

		return openFileError
	}
	defer file.Close()

	saveFileUpload, saveFileUploadError := uploadfile.SaveFileUpload(uploadfile.SaveFileUploadParams{
		UploadPath: resolvedFsPath,
		FileHeader: fileHeader,
		FileData:   file,
	})

	if saveFileUploadError != nil {
		log.Printf("Error while processing upload: %s", saveFileUploadError.Error())

		if saveFileUploadError == uploadfile.FileExistsError {
			echoContext.String(http.StatusConflict, "File already exists")

			return nil
		}

		echoContext.String(http.StatusBadRequest, "Failed to save file")

		return saveFileUploadError
	}

	log.Printf("File uploaded successfully to %s", saveFileUpload.FilePath)
	echoContext.String(http.StatusNoContent, "File uploaded successfully")

	return nil
}
