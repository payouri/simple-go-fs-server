package helpers

import (
	"errors"
	"path"
	"strings"
)

func IsPathLegal(basePath string, filePath string) (bool, error) {
	if !path.IsAbs(basePath) {
		return false, errors.New("base path is not absolute")
	}
	if !path.IsAbs(filePath) {
		return false, errors.New("file path is not absolute")
	}

	return strings.HasPrefix(filePath, basePath), nil
}
