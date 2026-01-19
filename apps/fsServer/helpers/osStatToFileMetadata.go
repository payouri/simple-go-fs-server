package helpers

import (
	"log"
	"os"
	"time"
)

type FileMetadata struct {
	Name         string `json:"name"`
	IsDir        bool   `json:"isDir"`
	BytesSize    int64  `json:"bytesSize"`
	LastModified string `json:"lastModified"`
	Permissions  string `json:"permissions"`
	MimeType     string `json:"mimeType"`
}

func OSStatToFileMetadata(location string, pathInfo os.FileInfo) FileMetadata {
	isDir := pathInfo.IsDir()
	if !isDir {
		mimeType, errorMimeType := GetAudioFileMimeType(location)
		if errorMimeType != nil {
			log.Printf("Error while getting mime type: %s", errorMimeType.Error())
		}

		return FileMetadata{
			Name:         pathInfo.Name(),
			IsDir:        isDir,
			BytesSize:    pathInfo.Size(),
			Permissions:  ConvertStringPermissionsToNumeric(pathInfo.Mode().String()),
			LastModified: pathInfo.ModTime().UTC().Format(time.RFC3339),
			MimeType:     mimeType,
		}
	}

	return FileMetadata{
		Name:         pathInfo.Name(),
		IsDir:        isDir,
		BytesSize:    pathInfo.Size(),
		Permissions:  ConvertStringPermissionsToNumeric(pathInfo.Mode().String()),
		LastModified: pathInfo.ModTime().UTC().Format(time.RFC3339),
	}
}
