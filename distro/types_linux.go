package distro

import (
	"fmt"
	"os"

	"github.com/jxsl13/osfacts/info"
)

type fileParseFunc func(dist distribution, filePath, fileContent string, osInfo *info.Os) error

type distribution struct {
	Name        string
	SearchNames []string
	Alias       string
	ParseFunc   fileParseFunc
}

func (o *distribution) search(fileContent string) error {
	_, err := mustContainOneOf(fileContent, unique(append(o.SearchNames, o.Name))...)
	return err
}

func (o *distribution) InfoName() string {
	if o.Alias != "" {
		return o.Alias
	}
	return o.Name
}

func (o *distribution) Parse(filePath, fileContent string) (*info.Os, error) {
	err := o.search(fileContent)
	if err != nil {
		return nil, err
	}

	parser := parseFallbackDistFile
	if o.ParseFunc != nil {
		parser = o.ParseFunc
	}

	osInfo := info.NewOs()

	err = parser(*o, filePath, fileContent, osInfo)
	if err != nil {
		return nil, err
	}
	return osInfo, nil
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
