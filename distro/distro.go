package distro

import (
	"errors"

	"github.com/jxsl13/osfacts/info"
)

var (
	ErrDetectionFailed = errors.New("failed to detect distribution")
)

func Detect() (*info.Os, error) {
	return detect()
}
