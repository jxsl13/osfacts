package distro

import (
	"github.com/jxsl13/osfacts/info"
)

func parseSlackwareDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	distName, err := mustContainOneOf(fileContent, dist.Name)
	if err != nil {
		return err
	}

	version, err := findSemanticVersion(fileContent)
	if err != nil {
		return err
	}

	return osInfo.Update(distName, version)
}
