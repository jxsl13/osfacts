package distro

import (
	"errors"
	"runtime"
)

var (
	ErrDetectionFailed = errors.New("failed to detect distribution")
)

func Detect() (*Info, error) {
	return detect()
}

// Info is the struct that contains the os information
// It i snot suposed to be used on its own but it is returned by
// the Detect() function.
type Info struct {
	// linux, windows etc
	Family string `json:"family"`
	// architecture: amd64, arm64, etc.
	Arch string `json:"arch"`
	// ubuntu, alpine, server (windows server)
	Distribution string `json:"distribution"`
	// 11.1.1
	Version string `json:"version"`
}

func newInfo() *Info {
	return &Info{
		Family: runtime.GOOS,
		Arch:   runtime.GOARCH,
	}
}

func (info *Info) update(distribution, version string) *Info {
	if len(version) > len(info.Version) && distribution != "" {
		info.Version = version
	}

	if distribution != "" {
		info.Distribution = distribution
	}
	return info
}
