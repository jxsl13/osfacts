package distro

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/jxsl13/osfacts/info"
)

func detectWithParams(filePath, fileContent string) (*info.Os, error) {
	distFilePaths := newDistMap()
	dists, found := distFilePaths[filePath]
	if !found {
		return nil, fmt.Errorf("%w: unknown path: %s", ErrDetectionFailed, filePath)
	}

	for _, dist := range dists {
		osInfo, err := dist.Parse(filePath, fileContent)
		if err != nil {
			continue
		}
		return osInfo, nil
	}

	return nil, ErrDetectionFailed
}

func Test_detectWithParams(t *testing.T) {
	tests := []struct {
		filePath    string
		fileContent string
		want        info.Os
		wantErr     bool
	}{
		{
			"/etc/os-release",
			`NAME="Alpine Linux"
ID=alpine
VERSION_ID=3.16.2
PRETTY_NAME="Alpine Linux v3.16"`,
			info.Os{Family: runtime.GOOS, Arch: runtime.GOARCH, Distribution: "Alpine Linux", Version: "3.16.2"},
			false,
		},
		{
			"/etc/os-release",
			`NAME="Ubuntu"
ID=ubuntu
VERSION="20.04.5 LTS (Focal Fossa)"
VERSION_ID="20.04"
ID_LIKE=debian
PRETTY_NAME="Ubuntu 20.04.5 LTS"`,
			info.Os{Family: runtime.GOOS, Arch: runtime.GOARCH, Distribution: "Ubuntu", Version: "20.04.5"},
			false,
		},
		{
			"/etc/centos-release",
			`CentOS release 6.10 (Final)`,
			info.Os{Family: runtime.GOOS, Arch: runtime.GOARCH, Distribution: "CentOS", Version: "6.10"},
			false,
		},
		{
			"/etc/os-release",
			`NAME="SLES"
VERSION="15-SP3"
VERSION_ID="15.3"
PRETTY_NAME="SUSE Linux Enterprise Server 15 SP3"
ID="sles"
ID_LIKE="suse"`,
			info.Os{Family: runtime.GOOS, Arch: runtime.GOARCH, Distribution: "SLES", Version: "15.3"},
			false,
		},
		{
			"/etc/redhat-release",
			`Red Hat Enterprise Linux Server release 6.10 (Santiago)`,
			info.Os{Family: runtime.GOOS, Arch: runtime.GOARCH, Distribution: "RedHat", Version: "6.10"},
			false,
		},
	}
	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%d: %s", idx+1, tt.want.Distribution), func(t *testing.T) {
			got, err := detectWithParams(tt.filePath, tt.fileContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectWithParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("detectWithParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
