package distro

import "github.com/jxsl13/osfacts/info"

func ParseAmazonDistFile(dist distribution, fileContent string) (*info.Os, error) {
	_, err := mustContainOneOf(fileContent, dist.Name)
	if err != nil {
		return nil, err
	}
	osInfo := info.NewOs()
	osInfo.Distribution = dist.Name

	switch dist.Path {
	case "/etc/os-release":
		version, err := findOsReleaseSemanticVersion(fileContent, "VERSION_ID")
		if err != nil {
			return nil, err
		}
		osInfo.Version = version
	default:
		version, err := findSemanticVersion(fileContent)
		if err != nil {
			return nil, err
		}
		osInfo.Version = version
	}

	return osInfo, nil
}
