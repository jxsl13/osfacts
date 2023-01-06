package distro

import (
	"fmt"
	"strings"

	"github.com/jxsl13/osfacts/info"
)

func ParseFallbackDistFile(dist distribution, fileContent string) (*info.Os, error) {
	osInfo := info.NewOs()
	osInfo.Distribution = dist.Name

	m, err := getOsReleaseMap(fileContent)
	if err == nil {
		name, err := getKey(m, "NAME")
		if err == nil {
			osInfo.Distribution = name
			version, err := findOsReleaseSemanticVersionInMap(m)
			if err == nil {
				osInfo.Version = version
				return osInfo, nil
			} else {
				version, err := findSemanticVersion(fileContent)
				if err != nil {
					return nil, err
				}
				osInfo.Version = version
				return osInfo, nil
			}
		} else {
			lines := strings.SplitN(fileContent, "\n\r", 2)
			if len(lines) == 0 {
				return nil, fmt.Errorf("%w: %s", ErrInvalidFileFormat, dist.Path)
			}
			tokens := strings.Split(lines[0], "\t\n\v\f\r ")
			name := tokens[0]
			if strings.Contains(name, "!§$%&/()=?\\") {
				return nil, fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, dist.Path)
			}
			osInfo.Distribution = name
			version, err := findOsReleaseSemanticVersionInMap(m)
			if err == nil {
				osInfo.Version = version
				return osInfo, nil
			} else {
				version, err := findSemanticVersion(fileContent)
				if err != nil {
					return nil, err
				}
				osInfo.Version = version
				return osInfo, nil
			}
		}
	} else {
		lines := strings.SplitN(fileContent, "\n\r", 2)
		if len(lines) == 0 {
			return nil, fmt.Errorf("%w: %s", ErrInvalidFileFormat, dist.Path)
		}
		tokens := strings.Split(lines[0], "\t\n\v\f\r ")
		name := tokens[0]
		if strings.Contains(name, "!§$%&/()=?\\") {
			return nil, fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, dist.Path)
		}
		osInfo.Distribution = name

		version, err := findSemanticVersion(fileContent)
		if err != nil {
			return nil, err
		}
		osInfo.Version = version
		return osInfo, nil
	}
}
