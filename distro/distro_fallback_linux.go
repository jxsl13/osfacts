package distro

import (
	"fmt"
	"strings"
)

func parseFallbackDistFile(dist distribution, filePath, fileContent string, osInfo *Info) error {

	m, err := getEnvMap(fileContent)
	if err == nil {
		name, err := getKey(m, "NAME")
		if err == nil {
			version, err := findEnvSemanticVersionInMap(m)
			if err == nil {
				osInfo.update(name, version)
				return nil
			} else {
				version, err := findSemanticVersionString(fileContent)
				if err != nil {
					return err
				}
				osInfo.update(name, version)
				return nil
			}
		} else {
			lines := strings.SplitN(fileContent, "\n\r", 2)
			if len(lines) == 0 {
				return fmt.Errorf("%w: %s", ErrInvalidFileFormat, filePath)
			}
			tokens := strings.Split(lines[0], "\t\n\v\f\r ")
			name := tokens[0]
			if strings.Contains(name, "!§$%&/()=?\\") {
				return fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, filePath)
			}

			version, err := findEnvSemanticVersionInMap(m)
			if err == nil {
				osInfo.update(name, version)
				return nil
			} else {
				version, err := findSemanticVersionString(fileContent)
				if err != nil {
					return err
				}
				osInfo.update(name, version)
				return nil
			}
		}
	} else {
		lines := strings.SplitN(fileContent, "\n\r", 2)
		if len(lines) == 0 {
			return fmt.Errorf("%w: %s", ErrInvalidFileFormat, filePath)
		}
		tokens := strings.Split(lines[0], "\t\n\v\f\r ")
		name := tokens[0]
		if strings.Contains(name, "!§$%&/()=?\\") {
			return fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, filePath)
		}

		version, err := findSemanticVersionString(fileContent)
		if err != nil {
			return err
		}
		osInfo.update(name, version)
		return nil
	}
}
