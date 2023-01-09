package distro

import (
	"github.com/jxsl13/osfacts/info"
)

func parseCentOSDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	envMap, err := getOsReleaseMap(fileContent)
	if err != nil {
		return err
	}

	nameValue, err := getKey(envMap, "NAME")
	if err != nil {
		return err
	}

	distName, err := mustContainOneOf(nameValue, "CentOS Linux", "CentOS Stream", "TencentOS")
	if err != nil {
		return err
	}

	switch distName {
	case "CentOS Linux":
		distName = "CentOS"
	}

	version, err := findOsReleaseSemanticVersionInMap(envMap, "VERSION_ID", "VERSION")
	if err != nil {
		return err
	}

	osInfo.Update(distName, version)
	return nil
}
