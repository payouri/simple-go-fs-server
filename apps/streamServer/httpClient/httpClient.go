package httpclient

import (
	"io"
	"media-stream-server/libs"
	"net/http"
	"os"
	"path"
	"time"
)

func getFsServerUrl() string {
	FS_SERVER_URL := os.Getenv("FS_SERVER_URL")
	if FS_SERVER_URL == "" {
		return "localhost:5000"
	}

	return FS_SERVER_URL
}
func getFsServerApiKey() string {
	FS_SERVER_API_KEY := os.Getenv("FS_SERVER_API_KEY")

	return FS_SERVER_API_KEY
}

type StreamForeignMediaParams struct {
	MediaId string `json:"media_id"`
}

type BuildHttpClientParams struct {
	fsServerUrl    string
	fsServerApiKey string
}
type HttpClientType struct {
	StreamForeignMedia func(params StreamForeignMediaParams) (string, error)
}

func buildStreamForeignMedia(dependencies BuildHttpClientParams) func(params StreamForeignMediaParams) (string, error) {
	var endpointURL = path.Join(dependencies.fsServerUrl, "/download")
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	return func(params StreamForeignMediaParams) (string, error) {
		mediaUrl := path.Join(endpointURL, params.MediaId)

		newRequest, errorRequest := http.NewRequest(http.MethodGet, mediaUrl, nil)
		if errorRequest != nil {
			return "", errorRequest
		}
		newRequest.Header.Set("Accept", "*/*")
		newRequest.Header.Set("Authorization", dependencies.fsServerApiKey)
		newRequest.Header.Set("User-Agent", "Stream-Server/1.0.0")

		response, errorResponse := client.Do(newRequest)
		if errorResponse != nil {
			return "", errorResponse
		}
		defer response.Body.Close()

		fileWriter, fileWriterError := libs.GetTmpFileWriter(params.MediaId)
		if fileWriterError != nil {
			return "", fileWriterError
		}

		body, errorBody := io.ReadAll(response.Body)
		if errorBody != nil {
			return "", errorBody
		}

		_, errorWrite := fileWriter.Write(body)
		if errorWrite != nil {
			return "", errorWrite
		}

		return "", nil
	}
}
func BuildHttpClient(dependencies BuildHttpClientParams) HttpClientType {
	stream := buildStreamForeignMedia(dependencies)

	return HttpClientType{
		StreamForeignMedia: stream,
	}
}

var HttpClient = BuildHttpClient(BuildHttpClientParams{
	fsServerUrl:    getFsServerUrl(),
	fsServerApiKey: getFsServerApiKey(),
})
