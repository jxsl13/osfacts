package distro

import "github.com/jxsl13/osfacts/info"

func parseAlpineDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	_, err := mustContainOneOf(fileContent, dist.Name)
	if err != nil {
		return err
	}
	osInfo.Distribution = dist.Name

	version, err := findOsReleaseSemanticVersion(fileContent)
	if err == nil {
		osInfo.Update(dist.Name, version)
		return nil
	}

	version, err = findSemanticVersion(fileContent)
	if err != nil {
		return err
	}

	osInfo.Update(dist.Name, version)
	return nil
}
