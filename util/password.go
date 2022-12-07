package util

import (
	"io"
	"os"
)

type PasswordProvider interface {
	GetPassword() (string, error)
}

type PasswordFilePath string

type filePasswordProvider struct {
	path PasswordFilePath
}

func NewFilePasswordProvider(path PasswordFilePath) PasswordProvider {
	return &filePasswordProvider{path: path}
}

func (f *filePasswordProvider) GetPassword() (string, error) {
	file, err := os.Open(string(f.path))
	if err != nil {
		return "", err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
