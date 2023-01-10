package distro

import (
	"fmt"

	"github.com/jxsl13/osfacts/info"
)

func parseCentOSDistFile(dist distribution, filePath, fileContent string, osInfo *info.Os) error {
	if filePath != "/etc/os-release" {
		return fmt.Errorf("invalid path for parser: %s", filePath)
	}

	envMap, err := getEnvMap(fileContent)
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

	version, err := findEnvSemanticVersionInMap(envMap, "VERSION_ID", "VERSION")
	if err != nil {
		return err
	}

	osInfo.Update(distName, version)
	return nil
}
