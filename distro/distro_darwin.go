package distro

import (
	"fmt"

	"github.com/jxsl13/osfacts/common"
	"github.com/jxsl13/osfacts/info"
)

func detect() (*info.Os, error) {
	cmd, err := common.NewCommand("/usr/bin/sw_vers", "-productVersion")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	output, err := cmd.OutputString()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	osInfo := info.NewOs()
	osInfo.Distribution = "macOS"

	version, err := findSemanticVersion(output)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	osInfo.Version = version
	return osInfo, nil
}
