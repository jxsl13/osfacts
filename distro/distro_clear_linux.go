package distro

import (
	"github.com/jxsl13/osfacts/info"
)

func parseClearLinuxDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	_, err := mustContainOneOf(fileContent, "Clear Linux")
	if err != nil {
		return err
	}

	version, err := findOsReleaseSemanticVersion(fileContent, "VERSION_ID")
	if err != nil {
		return err
	}

	return osInfo.Update(dist.Name, version)
}
