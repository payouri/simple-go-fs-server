package helpers

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
)

func ResolvePath(base string, path string) (string, error) {
	if path == "" {
		return "", errors.New("path to resolve is empty")
	}

	joinedPath := filepath.Join(base, path)
	joinedPath = strings.TrimPrefix(joinedPath, "/")
	joinedPath = strings.TrimSuffix(joinedPath, "/")

	if !fs.ValidPath(joinedPath) {
		return "", errors.New("path check failed")
	}

	localizeResult, localizeError := filepath.Localize(joinedPath)
	if localizeError != nil {
		return "", localizeError
	}

	return filepath.Clean("/" + localizeResult), nil
}
