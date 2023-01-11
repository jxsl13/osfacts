package distro

import (
	"github.com/jxsl13/osfacts/info"
)

type fileParseFunc func(dist distribution, filePath, fileContent string, osInfo *info.Os) error

type distribution struct {
	Name        string
	SearchNames []string
	Alias       string
	ParseFunc   fileParseFunc
}

type distPath struct {
	Path  string
	Dists []distribution
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
