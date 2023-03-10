package distro

type fileParseFunc func(dist distribution, filePath, fileContent string, osInfo *Info) error

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

func (o *distribution) Parse(filePath, fileContent string) (*Info, error) {
	err := o.search(fileContent)
	if err != nil {
		return nil, err
	}

	parser := parseFallbackDistFile
	if o.ParseFunc != nil {
		parser = o.ParseFunc
	}

	osInfo := newInfo()

	err = parser(*o, filePath, fileContent, osInfo)
	if err != nil {
		return nil, err
	}
	return osInfo, nil
}

func unique[T comparable](values []T) []T {
	m := make(map[T]struct{}, len(values))
	for _, v := range values {
		m[v] = struct{}{}
	}

	result := make([]T, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
