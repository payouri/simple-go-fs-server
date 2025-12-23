package httpClient

import (
	"fmt"
	"log"
	"media-stream-server/config"
	"net/http"
	"net/url"
	"os"
	"time"
)

func getFsServerUrl() string {
	FS_SERVER_URL := os.Getenv("FS_SERVER_URL")
	FS_SERVER_PORT := os.Getenv("FS_SERVER_PORT")

	if FS_SERVER_URL == "" {
		if FS_SERVER_PORT == "" {
			FS_SERVER_PORT = "5000"
		}
		return "http://localhost" + ":" + FS_SERVER_PORT
	}

	return FS_SERVER_URL
}
func getFsServerApiKey() string {
	FS_SERVER_API_KEY := os.Getenv("FS_SERVER_API_KEY")

	return FS_SERVER_API_KEY
}

type DownloadForeignMediaParams struct {
	MediaId string `json:"media_id"`
}

type BuildHttpClientParams struct {
	fsServerUrl    string
	fsServerApiKey string
}
type HttpClientType struct {
	DownloadForeignMedia func(params DownloadForeignMediaParams) (*http.Response, error)
}

func getResponseError(response *http.Response) error {
	if response.StatusCode >= 400 {
		return fmt.Errorf("Error while downloading foreign media, failed with status code %d", response.StatusCode)
	}

	return nil
}

func buildDownloadForeignMedia(dependencies BuildHttpClientParams) func(params DownloadForeignMediaParams) (*http.Response, error) {
	var endpointURL, endpointURLError = url.JoinPath(dependencies.fsServerUrl, "/download")
	if endpointURLError != nil {
		panic(endpointURLError)
	}

	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	return func(params DownloadForeignMediaParams) (*http.Response, error) {
		mediaUrl, joinUrlError := url.JoinPath(endpointURL, params.MediaId)
		if joinUrlError != nil {
			return nil, joinUrlError
		}

		newRequest, errorRequest := http.NewRequest(http.MethodGet, mediaUrl, nil)
		if errorRequest != nil {
			return nil, errorRequest
		}

		newRequest.Header.Set("Accept", "*/*")
		newRequest.Header.Set("Authorization", config.EnvConfig.FsServerApiKey)
		newRequest.Header.Set("User-Agent", "Stream-Server/1.0.0")

		response, errorResponse := client.Do(newRequest)
		if errorResponse != nil {
			return nil, errorResponse
		}
		hasError := getResponseError(response)
		if hasError != nil {
			return nil, hasError
		}

		log.Printf("Downloaded foreign media with id %s", params.MediaId)

		return response, nil
	}
}

func BuildHttpClient(dependencies BuildHttpClientParams) HttpClientType {
	download := buildDownloadForeignMedia(dependencies)

	return HttpClientType{
		DownloadForeignMedia: download,
	}
}

var HttpClient = BuildHttpClient(BuildHttpClientParams{
	fsServerUrl:    getFsServerUrl(),
	fsServerApiKey: getFsServerApiKey(),
})
