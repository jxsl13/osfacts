package distro

import (
	"fmt"

	"github.com/jxsl13/osfacts/info"
	"golang.org/x/sys/windows/registry"
)

func detect() (*info.Os, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("failed to access registry: %w", err)
	}
	defer k.Close()

	const (
		//distrubution
		productNameKey = "ProductName"

		// version
		displayVersionKey = "DisplayVersion"
		buildNumberKey    = "CurrentBuildNumber"
		patchLevelKey     = "UBR"
	)

	osInfo := info.NewOs()

	name, _, err := k.GetStringValue(productNameKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get value of %q: %w", productNameKey, err)
	}

	displayVersion, _, err := k.GetStringValue(displayVersionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get value of %q: %w", displayVersionKey, err)
	}

	buildNumber, _, err := k.GetStringValue(buildNumberKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get value of %q: %w", buildNumberKey, err)
	}

	patchLevel, _, err := k.GetIntegerValue(patchLevelKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get value of %q: %w", patchLevelKey, err)
	}

	osInfo.Distribution = name

	// same format as WinVer.exe
	osInfo.Version = fmt.Sprintf("Version %s (Build %s.%d)", displayVersion, buildNumber, patchLevel)
	return osInfo, nil
}
