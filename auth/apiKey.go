package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path"
	"simple-fs-web-service/helpers"
	"time"
)

const API_KEY_LENGTH = 32

func GenerateApiKey() (string, error) {
	// Create a byte slice of the desired length
	bytes := make([]byte, API_KEY_LENGTH)

	// Fill the slice with random bytes
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode the bytes as a base64 string
	apiKey := base64.URLEncoding.EncodeToString(bytes)

	// Trim to the desired length (base64 encoding increases length)
	return apiKey[:API_KEY_LENGTH], nil
}

func loadOfCreateFile(path string) (fs.File, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

var API_KEY_STORE_FILE_PATH = path.Join(helpers.GetBasePath(), "api_keys.json")

type ApiKeyData struct {
	ApiKey    string    `json:"apiKey"`
	ExpiresAt time.Time `json:"expiresAt"`
}
type ApiKeyStoreType struct {
	IsValid           func(string) bool
	GenerateApiKey    func() (string, error)
	HasOneExistingKey func() bool
}

func saveSnapshotToFile(path string, keyMap map[string]ApiKeyData) error {
	fileContent, marshallError := json.Marshal(keyMap)
	if marshallError != nil {
		return marshallError
	}

	writeError := os.WriteFile(API_KEY_STORE_FILE_PATH, fileContent, os.ModeExclusive)
	if writeError != nil {
		return writeError
	}

	return nil
}

func BuildApiKeyStore() ApiKeyStoreType {
	var keyMap = map[string]ApiKeyData{}
	file, err := loadOfCreateFile(API_KEY_STORE_FILE_PATH)
	if err != nil {
		log.Fatal("Error while creating file:", err)
	}
	stat, statErr := file.Stat()
	if statErr != nil {
		log.Fatal("Error while getting file stat:", statErr)
	}

	fileContents := make([]byte, stat.Size())
	_, readError := file.Read(fileContents)
	if readError != nil {
		log.Fatal("Error while reading file:", readError)
	}
	if len(fileContents) == 0 || stat.Size() == 0 {
		err := os.WriteFile(API_KEY_STORE_FILE_PATH, []byte("{}"), 0644)
		if err != nil {
			log.Fatal("Error while writing to file:", err)
		}
	} else {
		jsonParseError := json.Unmarshal(fileContents, &keyMap)
		if jsonParseError != nil {
			log.Fatal("Error while unmarshalling file:", err)
		}
	}
	defer file.Close()

	return ApiKeyStoreType{
		IsValid: func(s string) bool {
			if _, ok := keyMap[s]; !ok {
				return ok
			}

			return keyMap[s].ExpiresAt.After(time.Now())
		},
		HasOneExistingKey: func() bool {
			mapLen := len(keyMap)

			return mapLen > 0
		},
		GenerateApiKey: func() (string, error) {
			apiKey, err := GenerateApiKey()
			if err != nil {
				return "", err
			}
			keyMap[apiKey] = ApiKeyData{
				ApiKey:    apiKey,
				ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
			}

			saveError := saveSnapshotToFile(API_KEY_STORE_FILE_PATH, keyMap)
			if saveError != nil {
				return "", saveError
			}

			return apiKey, nil
		},
	}

}

var ApiKeyStore = BuildApiKeyStore()
