package helpers

import "os"

func GetBasePath() string {
	basePath := os.Getenv("BASE_PATH")
	homePath := os.Getenv("HOME")

	if basePath == "" {
		basePath = homePath
	}

	return basePath
}
