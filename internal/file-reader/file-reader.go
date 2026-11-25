package filereader

import (
	"errors"
	"os"
	"path"
	"strings"
)

var ErrFileFormatIsNotSupported = errors.New("file format is not supported")

func GetContext(fileName string) (string, error) {
	switch strings.ToLower(path.Ext(fileName)) {
	case ".txt", ".md":
		return getContextForText(fileName)
	}
	return "", ErrFileFormatIsNotSupported
}

func getContextForText(fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	return string(data), err
}
