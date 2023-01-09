package distro

import "github.com/jxsl13/osfacts/info"

func parseAmazonDistFile(dist distribution, fileContent string, osInfo *info.Os) error {
	_, err := mustContainOneOf(fileContent, dist.Name)
	if err != nil {
		return err
	}

	version := ""
	switch dist.Path {
	case "/etc/os-release":
		version, err = findOsReleaseSemanticVersion(fileContent, "VERSION_ID")
		if err != nil {
			return err
		}
		osInfo.Version = version
	default:
		version, err = findSemanticVersion(fileContent)
		if err != nil {
			return err
		}
		osInfo.Version = version
	}

	// version must be set at this point, otherwise an error should have occured
	osInfo.Update(dist.Name, version)
	return nil
}
