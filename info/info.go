package info

import "runtime"

type Os struct {
	// linux, windows etc
	Name string `json:"os"`
	// architecture: amd64, arm64, etc.
	Arch string `json:"arch"`
	// ubuntu, alpine, server (windows server)
	Distribution string `json:"distribution"`
	// 11.1.1
	Version string `json:"version"`
}

func NewOs() *Os {
	return &Os{
		Name: runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
