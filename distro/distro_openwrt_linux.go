package distro

import "github.com/jxsl13/osfacts/info"

func ParseOpenWrtDistFile(dist distribution, fileContent string) (*info.Os, error) {
	_, err := mustContainOneOf(fileContent, dist.Name)
	if err != nil {
		return nil, err
	}
	osInfo := info.NewOs()
	osInfo.Distribution = dist.Name

	version, err := findOsReleaseSemanticVersion(fileContent, "DISTRIB_RELEASE")
	if err == nil {
		osInfo.Version = version
		return osInfo, nil
	}

	version, err = findSemanticVersion(fileContent)
	if err != nil {
		return nil, err
	}
	osInfo.Version = version

	return osInfo, nil
}
