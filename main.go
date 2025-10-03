package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"simple-fs-web-service/routes/listDirectories"
	"simple-fs-web-service/validators"
	"strconv"
	"strings"
)

func getBasePath() string {
	basePath := os.Getenv("BASE_PATH")
	homePath := os.Getenv("HOME")

	if basePath == "" {
		basePath = homePath
	}

	return basePath
}

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

func handleUploadEndpoint(w http.ResponseWriter, r *http.Request, restSegments []string) {
	log.Printf("Upload endpoint called")
	path := strings.Join(restSegments, "/")
	uploadPath := filepath.Clean(filepath.Join(getBasePath(), path))

	parseFormError := r.ParseMultipartForm(10 * 1024 * 1024)
	if parseFormError != nil {
		log.Printf("Error while parsing form: %s", parseFormError.Error())
		http.Error(w, parseFormError.Error(), http.StatusBadRequest)
		return
	}
	file, fileHeader, fileError := r.FormFile("file")
	if fileError != nil {
		log.Printf("Error while getting file: %s", fileError.Error())
		http.Error(w, fileError.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileName := fileHeader.Filename
	filePath := filepath.Join(uploadPath, fileName)
	if _, err := os.Stat(filePath); err == nil {
		log.Printf("File %s already exists", filePath)
		http.Error(w, "file already exists", http.StatusBadRequest)
		return
	}
	uploadDirStat, err := os.Stat(uploadPath)
	if err != nil {
		log.Printf("Error while checking upload path: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !uploadDirStat.IsDir() {
		log.Printf("Upload path is not a directory")
		http.Error(w, "upload path is not a directory", http.StatusBadRequest)
		return
	}

	err = os.MkdirAll(uploadPath, 0755)
	if err != nil {
		log.Printf("Error while creating upload path: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Créer le fichier de destination
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error while creating destination file: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer dst.Close()

	// Copier le contenu du fichier source vers le fichier de destination
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error while copying file: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("File uploaded successfully to %s", filePath)

	w.WriteHeader(http.StatusOK)
}

func main() {
	port := os.Getenv("PORT")
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
		switch method {
		case http.MethodGet:
			switch firstSegment {
			case "list":
				restSegments := segments[1:]
				handleListEndpoint(w, r, restSegments)
			default:
				http.NotFound(w, r)
				log.Printf("Unknown request: %s", requestPath)
				return
			}
		case http.MethodPost:
			switch firstSegment {
			case "upload":
				restSegments := segments[1:]
				handleUploadEndpoint(w, r, restSegments)
			default:
				http.NotFound(w, r)
				log.Printf("Unknown request: %s", requestPath)
				return
			}
		}
	})

	log.Printf("Listening on port %s", port)

	serverAddress := ":" + port
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
