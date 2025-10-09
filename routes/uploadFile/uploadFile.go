package uploadfile

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type SaveFileUploadResponseSuccess struct {
	FilePath string
}

type SaveFileUploadParams struct {
	UploadPath string
	FileHeader *multipart.FileHeader
	FileData   multipart.File
}

func SaveFileUpload(params SaveFileUploadParams) (SaveFileUploadResponseSuccess, error) {
	uploadPath := params.UploadPath
	file := params.FileData
	fileHeader := params.FileHeader

	defer file.Close()
	fileName := fileHeader.Filename
	filePath := filepath.Join(uploadPath, fileName)
	if _, err := os.Stat(filePath); err == nil {
		log.Printf("File %s already exists", filePath)
		return SaveFileUploadResponseSuccess{}, errors.New("file already exists")
	}
	uploadDirStat, err := os.Stat(uploadPath)
	if err != nil {
		log.Printf("Error while checking upload path: %s", err.Error())
		return SaveFileUploadResponseSuccess{}, err
	}
	if !uploadDirStat.IsDir() {
		log.Printf("Upload path is not a directory")
		return SaveFileUploadResponseSuccess{}, err
	}

	err = os.MkdirAll(uploadPath, 0755)
	if err != nil {
		log.Printf("Error while creating upload path: %s", err.Error())
		return SaveFileUploadResponseSuccess{}, err
	}
	// Créer le fichier de destination
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error while creating destination file: %s", err.Error())
		return SaveFileUploadResponseSuccess{}, err
	}
	defer dst.Close()

	// Copier le contenu du fichier source vers le fichier de destination
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error while copying file: %s", err.Error())
		return SaveFileUploadResponseSuccess{}, err
	}

	return SaveFileUploadResponseSuccess{
		FilePath: filePath,
	}, nil
}
