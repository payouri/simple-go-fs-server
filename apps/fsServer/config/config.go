package config

import (
	"os"
)

type ConfigType struct {
	Port           string
	FsServerUrl    string
	FsServerApiKey string
}

func GetEnvConfig() ConfigType {
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5008"
	}

	FsServerUrl := os.Getenv("FS_SERVER_URL")
	FsServerPort := os.Getenv("FS_SERVER_PORT")
	if FsServerUrl == "" {
		if FsServerPort == "" {
			FsServerPort = "5008"
		}

		FsServerUrl = "localhost:" + FsServerPort
	}
	FsServerApiKey := os.Getenv("FS_SERVER_API_KEY")

	return ConfigType{
		Port:           Port,
		FsServerUrl:    FsServerUrl,
		FsServerApiKey: FsServerApiKey,
	}
}

var EnvConfig = GetEnvConfig()
