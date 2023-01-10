package distro

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jxsl13/osfacts/info"
)

func parseSuseReleaseDistFile(dist distribution, filePath, fileContent string, osInfo *info.Os) error {

	if filePath != "/etc/SuSE-release" {
		return errors.New("invalid SUSE release path")
	}
	distName, version := "", ""

	lines := strings.Split(fileContent, "\n")
	if len(lines) < 3 {
		// name, version, patch level
		return fmt.Errorf("%w: %s", ErrInvalidFileFormat, filePath)
	}
	lcFileContent := strings.ToLower(fileContent)

	if strings.Contains(lcFileContent, "open") {
		distLine := lines[0]
		tokens := strings.SplitN(distLine, " ", 2)
		if len(tokens) != 2 {
			return fmt.Errorf("%w: unexpected first line: %s", ErrInvalidFileFormat, filePath)
		}
		distName = tokens[0]
	} else if strings.Contains(lcFileContent, "enterprise") {

		if strings.Contains(fileContent, "Server") {
			distName = "SLES"
		} else if strings.Contains(fileContent, "Desktop") {
			distName = "SLED"
		}
	}

	realPath, err := os.Readlink("/etc/products.d/baseproduct")
	if err == nil && strings.HasSuffix(realPath, "SLES_SAP.prod") {
		distName = "SLES_SAP"
	}

	versionLines := lines[1:]
	re := regexp.MustCompile(`\s*=\s*`)

	for idx, line := range versionLines {
		versionLines[idx] = strings.TrimSpace(re.ReplaceAllString(line, "="))
	}
	versionData := strings.Join(versionLines, "\n")

	envMap, err := getEnvMap(versionData)
	if err != nil {
		// look inside of the whole file for a semantic version
		// in case we could not parse key/value pairs
		version, err := findSemanticVersion(fileContent)
		if err != nil {
			return err
		}
		osInfo.Update(distName, version)
		return nil
	}

	version, foundVersion := envMap["VERSION"]
	patchLevel, foundPatchLevel := envMap["PATCHLEVEL"]
	if !foundVersion || !foundPatchLevel {
		// look inside of the whole file for a semantiv version
		// in case we could not find out expected key names
		version, err := findSemanticVersion(fileContent)
		if err != nil {
			return err
		}
		osInfo.Update(distName, version)
		return nil
	}

	osInfo.Update(distName, fmt.Sprintf("%s.%s", version, patchLevel))
	return nil
}
