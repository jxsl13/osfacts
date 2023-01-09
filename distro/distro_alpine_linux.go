package distro

import "github.com/jxsl13/osfacts/info"

func parseAlpineDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	_, err := mustContainOneOf(fileContent, dist.Name)
	if err != nil {
		return err
	}

	version, err := findOsReleaseSemanticVersion(fileContent)
	if err == nil {
		return osInfo.Update(dist.Name, version)
	}

	version, err = findSemanticVersion(fileContent)
	if err != nil {
		return err
	}

	return osInfo.Update(dist.Name, version)
}
