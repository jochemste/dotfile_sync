package utils

import (
	"errors"
	"os"
	//"github.com/pelletier/go-toml/v2"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func GetToml(path string) (string, error) {
	if !FileExists(path) {
		return "", errors.New("File does not exist")
	}

	return "", nil
}
