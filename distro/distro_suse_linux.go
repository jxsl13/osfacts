package distro

import (
	"fmt"
	"strings"

	"github.com/jxsl13/osfacts/info"
)

func ParseSUSEDistFile(dist distribution, fileContent string) (*info.Os, error) {
	_, err := mustContainOneOf(fileContent, dist.Name, strings.ToLower(dist.Name))
	if err != nil {
		return nil, err
	}
	osInfo := info.NewOs()
	osInfo.Distribution = dist.Name

	switch dist.Path {
	case "/etc/os-release":
		m, err := getOsReleaseMap(fileContent)
		if err != nil {
			return nil, err
		}
		name, err := getKey(m, "NAME")
		if err != nil {
			return nil, err
		}
		osInfo.Distribution = name

		version, err := findOsReleaseSemanticVersionInMap(m, "VERSION_ID")
		if err != nil {
			return nil, err
		}
		osInfo.Version = version

	case "/etc/SuSE-release":
		lines := strings.Split(fileContent, "\n")
		if len(lines) == 0 {
			return nil, fmt.Errorf("%w: %s", ErrInvalidFileFormat, dist.Path)
		}
		lcFileContent := strings.ToLower(fileContent)

		if strings.Contains(lcFileContent, "open") {
			distLine := lines[0]
			tokens := strings.SplitN(distLine, " ", 2)
			if len(tokens) != 2 {
				return nil, fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, dist.Path)
			}
			osInfo.Distribution = tokens[0]
		} else if strings.Contains(lcFileContent, "enterprise") {

			if strings.Contains(fileContent, "Server") {
				osInfo.Distribution = "SLES"
			} else if strings.Contains(fileContent, "Desktop") {
				osInfo.Distribution = "SLED"
			}
		}

		version, err := findSemanticVersion(fileContent)
		if err != nil {
			return nil, err
		}
		osInfo.Version = version
	}

	return osInfo, nil
}
