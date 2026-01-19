package routes

import (
	routesHelpers "media-stream-server/routers/helpers"
	streamAudioFileController "media-stream-server/routes/streamAudioFile"
	"net/http"

	"github.com/labstack/echo/v4"
)

var STREAM_AUDIO_FILE_ROUTE_PATH = routesHelpers.GetRoutePath("stream")

func StreamAudioFileRoute(echoContext echo.Context) error {
	streamResult, streamError := streamAudioFileController.StreamAudioFile(streamAudioFileController.StreamAudioFileParams{
		MediaId: "Downloads/glitch-177348.mp3",
	})
	if streamError != nil {
		echoContext.String(http.StatusInternalServerError, streamError.Error())

		return streamError
	}
	if streamResult == nil {
		echoContext.String(http.StatusInternalServerError, "Stream result is nil")

		return nil
	}

	return echoContext.Stream(http.StatusAccepted, "application/octet-stream", streamResult)
}
