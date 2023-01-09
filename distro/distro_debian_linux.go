package distro

import "github.com/jxsl13/osfacts/info"

func parseDebianDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	distName, err := mustContainOneOf(fileContent, dist.Name, "Raspbian", "Ubuntu", "Devuan", "Cumulus", "Mint", "Deepin", "deepin", "UOS", "Uos", "uos")
	if err != nil {
		return err
	}
	defaultVersionSearch := true
	distVersion := ""

	switch distName {
	case "Debian", "Raspbian":
		defaultVersionSearch = false

		versionContent, err := getFileContent("/etc/debian_version")
		if err != nil {
			return err
		}
		distVersion, err = findSemanticVersion(versionContent)
		if err != nil {
			return err
		}

	case "Cumulus":
		distName = "Cumulus Linux"
	case "Mint":
		distName = "Linux Mint"
	case "UOS", "Uos", "uos":
		distName = "Uos"
	case "Deepin", "deepin":
		distName = "Deepin"
	default:
		//Ubuntu, Devuan
	}

	if defaultVersionSearch {
		distVersion, err = findOsReleaseSemanticVersion(fileContent, "VERSION_ID", "VERSION")
		if err != nil {
			return err
		}
	}

	return osInfo.Update(distName, distVersion)
}
