package helpers

import (
	"log"
	"net/http"
	"os"
	"path"
)

func detectMimeTypeFromExtension(extension string) string {
	switch extension {
	case ".mp3":
		return "audio/mpeg"
	case ".mp4":
		return "audio/mp4"
	case ".ogg":
		return "audio/ogg"
	case ".wav":
		return "audio/wav"
	case ".aac":
		return "audio/aac"
	case ".webm":
		return "audio/webm"
	default:
		return ""
	}
}

func GetAudioFileMimeType(filePath string) (string, error) {
	if !path.IsAbs(filePath) {
		return "", os.ErrInvalid
	}

	file, fileError := os.OpenFile(filePath, os.O_RDONLY, 0)
	if fileError != nil {
		return "", fileError
	}

	firstBytes := make([]byte, 512)
	read, readError := file.Read(firstBytes)
	if readError != nil {
		return "", readError
	}
	log.Printf("Read %d bytes", read)
	closeError := file.Close()
	if closeError != nil {
		return "", closeError
	}

	mimeType := http.DetectContentType(firstBytes)

	if mimeType == "application/octet-stream" {
		detectedMimeType := detectMimeTypeFromExtension(path.Ext(filePath))

		if detectedMimeType != "" {
			return detectedMimeType, nil
		}
	}

	return mimeType, nil
}
