package distro

import "github.com/jxsl13/osfacts/info"

func parseOpenWrtDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	distName, err := mustContainOneOf(fileContent, dist.Name)
	if err != nil {
		return err
	}

	version, err := findOsReleaseSemanticVersion(fileContent, "DISTRIB_RELEASE")
	if err == nil {
		osInfo.Update(distName, version)
		return nil
	}

	version, err = findSemanticVersion(fileContent)
	if err != nil {
		return err
	}

	osInfo.Update(distName, version)
	return nil
}
