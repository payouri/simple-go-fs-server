package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"simple-fs-web-service/auth"
	"simple-fs-web-service/helpers"
	listDirectories "simple-fs-web-service/routes/listDirectories"
	uploadfile "simple-fs-web-service/routes/uploadFile"
	"simple-fs-web-service/validators"
	"strconv"
	"strings"
)

func handleListEndpoint(w http.ResponseWriter, r *http.Request, restSegments []string) {
	limitInt, limitErr := validators.ParseLimit(r.URL.Query().Get("limit"))
	offsetInt, offsetErr := validators.ParseOffset(r.URL.Query().Get("offset"))
	parseSortFieldResult, parseSortFieldErr := validators.ParseSortField(r.URL.Query().Get("sortField"))
	parseSortOrderResult, parseSortOrderErr := validators.ParseSortOrder(r.URL.Query().Get("sortOrder"))
	if limitErr != nil {
		log.Printf("Limit is not a valid integer: %s", limitErr.Error())
		http.Error(w, "limit is not a valid integer", http.StatusBadRequest)
		return
	}
	if offsetErr != nil {
		log.Printf("Offset is not a valid integer: %s", offsetErr.Error())
		http.Error(w, "offset is not a valid integer", http.StatusBadRequest)
		return
	}
	if parseSortFieldErr != nil {
		log.Printf("Sort field is not a valid: %s", parseSortFieldErr.Error())
		http.Error(w, "sort field is not a valid", http.StatusBadRequest)
		return
	}
	if parseSortOrderErr != nil {
		log.Printf("Sort order is not a valid: %s", parseSortOrderErr.Error())
		http.Error(w, "sort order is not a valid", http.StatusBadRequest)
		return
	}
	listDirectoryResult, listDirectoryError := listDirectories.ListDirectory(listDirectories.ListDirectoryParams{
		Path:      strings.Join(restSegments, "/"),
		Limit:     limitInt,
		Offset:    offsetInt,
		SortField: parseSortFieldResult,
		SortOrder: parseSortOrderResult,
	})
	if listDirectoryError != nil {
		http.Error(w, listDirectoryError.Error(), http.StatusBadRequest)
		log.Printf("Error while listing directory: %s", listDirectoryError.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"files": listDirectoryResult.Files, "total": listDirectoryResult.Total, "maxPage": listDirectoryResult.MaxPage, "page": listDirectoryResult.Page}
	json.NewEncoder(w).Encode(response)
}

func handleDownloadEndpoint(w http.ResponseWriter, r *http.Request, restSegments []string) {
	path := strings.Join(restSegments, "/")
	downloadPath := filepath.Clean(filepath.Join(helpers.GetBasePath(), path))
	file, fileError := os.Open(downloadPath)
	if fileError != nil {
		log.Printf("Error while opening file: %s", fileError.Error())
		http.Error(w, fileError.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(downloadPath)))
	io.Copy(w, file)
}

const MAX_UPLOAD_SIZE = 100 * 1024 * 1024

func handleUploadEndpoint(w http.ResponseWriter, r *http.Request, restSegments []string) {
	log.Printf("Upload endpoint called")
	path := strings.Join(restSegments, "/")
	uploadPath := filepath.Clean(filepath.Join(helpers.GetBasePath(), path))
	file, fileHeader, fileError := r.FormFile("file")

	if fileError != nil {
		log.Printf("Error while getting file: %s", fileError.Error())
		http.Error(w, fileError.Error(), http.StatusBadRequest)
		return
	}

	if fileHeader.Size > MAX_UPLOAD_SIZE {
		log.Printf("File too big: %vMB", fileHeader.Size/1024/1024)
		http.Error(w, "file too big", http.StatusRequestEntityTooLarge)
		return
	}

	saveFileUpload, saveFileUploadError := uploadfile.SaveFileUpload(uploadfile.SaveFileUploadParams{
		UploadPath: uploadPath,
		FileHeader: fileHeader,
		FileData:   file,
	})

	if saveFileUploadError != nil {
		log.Printf("Error while processing upload: %s", saveFileUploadError.Error())
		http.Error(w, saveFileUploadError.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("File uploaded successfully to %s", saveFileUpload.FilePath)

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	port := os.Getenv("PORT")

	log.Printf("exist, %v", auth.ApiKeyStore.HasOneExistingKey())
	if !auth.ApiKeyStore.HasOneExistingKey() {
		_, generateKeyError := auth.ApiKeyStore.GenerateApiKey()

		if generateKeyError != nil {
			log.Fatal(generateKeyError)
		}
	}

	if port == "" {
		port = "5000"
	}
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Le port %s n'est pas un port valide.", port)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		method := r.Method
		log.Printf("Request %v", r.URL)
		without, _ := strings.CutPrefix(requestPath, "/")
		segments := strings.Split(without, "/")
		firstSegment := segments[0]
		restSegments := segments[1:]
		switch method {
		case http.MethodGet:
			switch firstSegment {
			case "list":
				auth.AuthRequest(handleListEndpoint)(w, r, restSegments)
				return
			case "download":
				auth.AuthRequest(handleDownloadEndpoint)(w, r, restSegments)
				return
			default:
				http.NotFound(w, r)
				log.Printf("Unknown request: %s", requestPath)
				return
			}
		case http.MethodPost:
			switch firstSegment {
			case "upload":
				auth.AuthRequest(handleUploadEndpoint)(w, r, restSegments)
				return
			default:
				http.NotFound(w, r)
				log.Printf("Unknown request: %s", requestPath)
				return
			}
		default:
			http.NotFound(w, r)
			log.Printf("Unknown request: %s", requestPath)
			return
		}
	})

	log.Printf("Listening on port %s", port)

	serverAddress := ":" + port
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
