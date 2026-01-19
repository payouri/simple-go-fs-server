package getFileMetadataRoute

import (
	"errors"
	"os"
	"simple-fs-web-service/helpers"
)

type GetFileMetadataParams struct {
	FilePath string `json:"file_path"`
}

type GetFileMetadataResponseSuccess = helpers.FileMetadata

func GetFileMetadata(params GetFileMetadataParams) (GetFileMetadataResponseSuccess, error) {
	absoluteFilePath, resolvePathError := helpers.ResolvePath(helpers.GetBasePath(), params.FilePath)
	if resolvePathError != nil {
		return GetFileMetadataResponseSuccess{}, resolvePathError
	}
	isPathLegalResult, isPathLegalError := helpers.IsPathLegal(helpers.GetBasePath(), absoluteFilePath)
	if isPathLegalError != nil {
		return GetFileMetadataResponseSuccess{}, isPathLegalError
	}
	if !isPathLegalResult {
		return GetFileMetadataResponseSuccess{}, errors.New("path is not a valid path")
	}

	pathInfo, pathInfoError := os.Stat(absoluteFilePath)
	if pathInfoError != nil {
		return GetFileMetadataResponseSuccess{}, pathInfoError
	}
	if pathInfo.IsDir() {
		return GetFileMetadataResponseSuccess{}, errors.New("path is a directory")
	}

	return helpers.OSStatToFileMetadata(absoluteFilePath, pathInfo), nil
}
