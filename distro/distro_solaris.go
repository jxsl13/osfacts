package distro

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/sys/unix"
)

func detect() (*Info, error) {
	fileContent, err := getFileContent("/etc/release")
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(fileContent, "\n", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, ErrInvalidFileFormat)
	}

	firstLine := parts[0]

	utsName := unix.Utsname{}
	distName := ""

	err = unix.Uname(&utsName)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	uname := strings.Join([]string{string(utsName.Nodename[:]), string(utsName.Release[:]), string(utsName.Version[:])}, " ")
	distName, err = mustContainOneOf(uname, "NexentaOS")
	if err != nil {
		distName, err = mustContainOneOf(firstLine, "Oracle Solaris", "SmartOS", "OpenIndiana", "OmniOS", "Solaris")
		if err != nil {
			return nil, fmt.Errorf("%w: unknown distribution name", ErrDetectionFailed)
		}
	}

	// rename
	switch distName {
	case "Solaris", "Oracle Solaris":
		distName = "Solaris"
	case "NexentaOS":
		distName = "Nexenta"
	}

	// version
	version := ""

	switch distName {
	case "SmartOS":
		productData, err := getFileContent("/etc/product")
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
		}

		re := regexp.MustCompile(`\s*:\s*`)
		versionLines := make([]string, 0, 4)
		for _, line := range strings.Split(productData, "\n") {
			if strings.Contains(line, ": ") {
				kvLine := re.ReplaceAllString(strings.TrimSpace(line), "=")
				versionLines = append(versionLines, kvLine)
			}
		}

		envMap, err := getEnvMap(strings.Join(versionLines, "\n"))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
		}

		versionString, found := envMap["Image"]
		if !found {
			return nil, fmt.Errorf("%w: version not found", ErrDetectionFailed)
		}

		version, err = findSemanticVersionString(versionString)
		if err != nil {
			return nil, fmt.Errorf("%w: version not found: %v", ErrDetectionFailed, err)
		}
	default:
		// "SmartOS", "OpenIndiana", "OmniOS", "Solaris", "Nexenta"

		version, err = findSemanticVersionString(firstLine)
		if err != nil {
			return nil, fmt.Errorf("%w: version not found: %v", ErrDetectionFailed, err)
		}
	}

	return newInfo().update(distName, version), nil
}
