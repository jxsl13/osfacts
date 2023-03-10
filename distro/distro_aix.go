package distro

import (
	"fmt"
	"strings"

	"github.com/jxsl13/osfacts/internal"
)

func detect() (*Info, error) {
	cmd, err := internal.NewCommand("/usr/bin/oslevel")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	output, err := cmd.OutputString()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	osInfo := newInfo()

	osInfo.Version = strings.TrimSpace(output)
	return osInfo, nil
}
