package routes

import (
	"errors"
	"media-stream-server/httpClient"
	routesHelpers "media-stream-server/routers/helpers"
	streamAudioFileController "media-stream-server/routes/streamAudioFile"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var STREAM_AUDIO_FILE_ROUTE_PATH = routesHelpers.GetRoutePath("stream/*")

func validateMediaMimeType(mimeType string) error {
	if !strings.HasPrefix(mimeType, "audio/") && !strings.HasPrefix(mimeType, "video/") && !strings.HasPrefix(mimeType, "application/octet-stream") {
		return errors.New("Unsupported stream mime type")
	}

	return nil
}

func StreamAudioFileRoute(echoContext echo.Context) error {
	mediaId := routesHelpers.GetFSFilePath(echoContext.Request().URL.Path, STREAM_AUDIO_FILE_ROUTE_PATH)
	getForeignMediaMetadataResult, getForeignMediaMetadataError := httpClient.HttpClient.GetForeignMediaMetadata(httpClient.GetForeignMediaMetadataParams{
		MediaId: mediaId,
	})
	if getForeignMediaMetadataError != nil {
		// echoContext.Logger().Error(getForeignMediaMetadataError.Error())
		echoContext.String(http.StatusBadRequest, "Failed to retrieve stream metadata")

		return getForeignMediaMetadataError
	}
	if getForeignMediaMetadataResult.IsDir {
		echoContext.String(http.StatusBadRequest, "Invalid media id")

		return nil
	}
	if err := validateMediaMimeType(getForeignMediaMetadataResult.MimeType); err != nil {
		echoContext.Logger().Error(err.Error())
		echoContext.String(http.StatusBadRequest, "Unsupported media mime type")

		return err
	}
	echoContext.Response().Header().Set("cache", "no-cache")
	streamResult, streamError := streamAudioFileController.StreamAudioFile(streamAudioFileController.StreamAudioFileParams{
		MediaId: mediaId,
	})
	if streamError != nil {
		echoContext.String(http.StatusInternalServerError, streamError.Error())

		return streamError
	}
	if streamResult == nil {
		echoContext.String(http.StatusInternalServerError, "Stream result is nil")

		return nil
	}

	return echoContext.Stream(http.StatusPartialContent, "audio/webm", streamResult)
}
