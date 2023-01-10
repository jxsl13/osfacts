package distro

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jxsl13/osfacts/info"
)

func parseSuseReleaseDistFile(dist distribution, filePath, fileContent string, osInfo *info.Os) error {

	if filePath != "/etc/SuSE-release" {
		return errors.New("invalid SUSE release path")
	}
	distName, version := "", ""

	lines := strings.Split(fileContent, "\n")
	if len(lines) == 0 {
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

	version, err := findSemanticVersion(fileContent)
	if err != nil {
		return err
	}

	realPath, err := os.Readlink("/etc/products.d/baseproduct")
	if err == nil && strings.HasSuffix(realPath, "SLES_SAP.prod") {
		distName = "SLES_SAP"
	}

	osInfo.Update(distName, version)
	return nil
}
