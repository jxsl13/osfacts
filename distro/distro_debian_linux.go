package distro

import "github.com/jxsl13/osfacts/info"

func ParseDebianDistFile(dist distribution, fileContent string) (*info.Os, error) {
	distName, err := mustContainOneOf(fileContent, dist.Name, "Raspbian", "Ubuntu", "Devuan", "Cumulus", "Mint", "Deepin", "deepin", "UOS", "Uos", "uos")
	if err != nil {
		return nil, err
	}
	osInfo := info.NewOs()
	defaultVersionSearch := true

	switch distName {
	case "Debian", "Raspbian":
		defaultVersionSearch = false
		versionContent, err := getFileContent("/etc/debian_version")
		if err != nil {
			return nil, err
		}
		version, err := findSemanticVersion(versionContent)
		if err != nil {
			return nil, err
		}
		osInfo.Version = version

	case "Cumulus":
		osInfo.Distribution = "Cumulus Linux"
	case "Mint":
		osInfo.Distribution = "Linux Mint"
	case "UOS", "Uos", "uos":
		osInfo.Distribution = "Uos"
	case "Deepin", "deepin":
		osInfo.Distribution = "Deepin"
	default:
		osInfo.Distribution = distName
		//Ubuntu, Devuan
	}

	if defaultVersionSearch {
		version, err := findOsReleaseSemanticVersion(fileContent, "VERSION_ID", "VERSION")
		if err != nil {
			return nil, err
		}
		osInfo.Version = version
	}

	return osInfo, nil
}
