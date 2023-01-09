package distro

import (
	"fmt"
	"strings"

	"github.com/jxsl13/osfacts/common"
	"github.com/jxsl13/osfacts/info"
)

func detect() (*info.Os, error) {
	cmd, err := common.NewCommand("/usr/bin/oslevel")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	output, err := cmd.OutputString()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	osInfo := info.NewOs()

	osInfo.Version = strings.TrimSpace(output)
	return osInfo, nil
}
