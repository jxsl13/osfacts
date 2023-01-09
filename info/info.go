package info

import (
	"errors"
	"runtime"
)

type Os struct {
	// linux, windows etc
	Family string `json:"family"`
	// architecture: amd64, arm64, etc.
	Arch string `json:"arch"`
	// ubuntu, alpine, server (windows server)
	Distribution string `json:"distribution"`
	// 11.1.1
	Version string `json:"version"`
}

func NewOs() *Os {
	return &Os{
		Family: runtime.GOOS,
		Arch:   runtime.GOARCH,
	}
}

func (info *Os) Update(distribution, version string) error {
	if len(version) > len(info.Version) && distribution != "" {
		info.Distribution = distribution
		info.Version = version
		return nil
	}
	return errors.New("not updated")
}
