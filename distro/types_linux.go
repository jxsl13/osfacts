package distro

import (
	"fmt"
	"os"

	"github.com/jxsl13/osfacts/info"
)

type fileParseFunc func(dist distribution, fileContent string, osInfo *info.Os) error

type distribution struct {
	Name string
	Path string
}

func (o *distribution) Content() (string, error) {
	return getFileContent(o.Path)
}

func getFileContent(path string) (string, error) {
	found, err := existsWithSize(path, false)
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	if !found {
		return "", fmt.Errorf("%s: not found or empty", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	return string(data), nil
}
