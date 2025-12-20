package libs

import (
	"io"
	"os"
	"path"
)

func GetTmpFilePath(fileName string) string {
	tmpFilePath := path.Join(os.TempDir(), fileName)
	return tmpFilePath
}

func GetTmpFile(fileName string) (*os.File, error) {
	tmpFilePath := GetTmpFilePath(fileName)
	file, errorFile := os.Create(tmpFilePath)
	if errorFile != nil {
		return nil, errorFile
	}
	return file, nil
}

func GetTmpFileReader(fileName string) (io.ReadCloser, error) {
	tmpFilePath := GetTmpFilePath(fileName)
	file, errorFile := os.Open(tmpFilePath)
	if errorFile != nil {
		return nil, errorFile
	}
	return file, nil
}

func GetTmpFileWriter(fileName string) (io.WriteCloser, error) {
	tmpFilePath := GetTmpFilePath(fileName)
	file, errorFile := os.OpenFile(tmpFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if errorFile != nil {
		return nil, errorFile
	}
	return file, nil
}
