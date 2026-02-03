package streamAudioFileController

import (
	"io"
	"log"
	"media-stream-server/httpClient"
	"media-stream-server/libs"
)

type StreamAudioFileParams struct {
	MediaId string `json:"media_id"`
}

func StreamAudioFile(params StreamAudioFileParams) (*io.PipeReader, error) {
	response, errorResponse := httpClient.HttpClient.DownloadForeignMedia(httpClient.DownloadForeignMediaParams{MediaId: params.MediaId})
	if errorResponse != nil {
		return nil, errorResponse
	}

	transcodeResult, transcodeError := libs.TranscodingClient.InMemoryTranscodeAudio(libs.InMemoryTranscodeAudioParams{
		InputStream:   response.Body,
		ContentLength: response.ContentLength,
		OutputFormat:  libs.MP4,
	})
	if transcodeError != nil {
		log.Printf("transcodeError %v", transcodeError)

		return nil, transcodeError
	}

	return transcodeResult, nil

}
