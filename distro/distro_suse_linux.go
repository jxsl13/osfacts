package distro

import (
	"fmt"
	"strings"

	"github.com/jxsl13/osfacts/info"
)

func parseSUSEDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	_, err := mustContainOneOf(fileContent, dist.Name, strings.ToLower(dist.Name))
	if err != nil {
		return err
	}

	switch dist.Path {
	case "/etc/os-release":
		m, err := getOsReleaseMap(fileContent)
		if err != nil {
			return err
		}
		distName, err := getKey(m, "NAME")
		if err != nil {
			return err
		}

		version, err := findOsReleaseSemanticVersionInMap(m, "VERSION_ID")
		if err != nil {
			return err
		}

		osInfo.Update(distName, version)

	case "/etc/SuSE-release":
		distName, version := "", ""

		lines := strings.Split(fileContent, "\n")
		if len(lines) == 0 {
			return fmt.Errorf("%w: %s", ErrInvalidFileFormat, dist.Path)
		}
		lcFileContent := strings.ToLower(fileContent)

		if strings.Contains(lcFileContent, "open") {
			distLine := lines[0]
			tokens := strings.SplitN(distLine, " ", 2)
			if len(tokens) != 2 {
				return fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, dist.Path)
			}
			distName = tokens[0]
		} else if strings.Contains(lcFileContent, "enterprise") {

			if strings.Contains(fileContent, "Server") {
				distName = "SLES"
			} else if strings.Contains(fileContent, "Desktop") {
				distName = "SLED"
			}
		}

		version, err := findSemanticVersion(fileContent)
		if err != nil {
			return err
		}

		osInfo.Update(distName, version)
	}

	return nil
}
