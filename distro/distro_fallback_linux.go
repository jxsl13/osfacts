package distro

import (
	"fmt"
	"strings"

	"github.com/jxsl13/osfacts/info"
)

func parseFallbackDistFile(dist distribution, fileContent string, osInfo *info.Os) error {

	m, err := getOsReleaseMap(fileContent)
	if err == nil {
		name, err := getKey(m, "NAME")
		if err == nil {
			version, err := findOsReleaseSemanticVersionInMap(m)
			if err == nil {
				osInfo.Update(name, version)
				return nil
			} else {
				version, err := findSemanticVersion(fileContent)
				if err != nil {
					return err
				}
				osInfo.Update(name, version)
				return nil
			}
		} else {
			lines := strings.SplitN(fileContent, "\n\r", 2)
			if len(lines) == 0 {
				return fmt.Errorf("%w: %s", ErrInvalidFileFormat, dist.Path)
			}
			tokens := strings.Split(lines[0], "\t\n\v\f\r ")
			name := tokens[0]
			if strings.Contains(name, "!§$%&/()=?\\") {
				return fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, dist.Path)
			}

			version, err := findOsReleaseSemanticVersionInMap(m)
			if err == nil {
				osInfo.Update(name, version)
				return nil
			} else {
				version, err := findSemanticVersion(fileContent)
				if err != nil {
					return err
				}
				osInfo.Update(name, version)
				return nil
			}
		}
	} else {
		lines := strings.SplitN(fileContent, "\n\r", 2)
		if len(lines) == 0 {
			return fmt.Errorf("%w: %s", ErrInvalidFileFormat, dist.Path)
		}
		tokens := strings.Split(lines[0], "\t\n\v\f\r ")
		name := tokens[0]
		if strings.Contains(name, "!§$%&/()=?\\") {
			return fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, dist.Path)
		}

		version, err := findSemanticVersion(fileContent)
		if err != nil {
			return err
		}
		osInfo.Update(name, version)
		return nil
	}
}
