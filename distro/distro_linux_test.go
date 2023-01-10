package distro

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/jxsl13/osfacts/info"
	"github.com/stretchr/testify/assert"
)

func detectWithParams(filePath, fileContent string) (*info.Os, error) {
	distPaths := newPaths()

	found := false
	var dp distPath
	for _, distPath := range distPaths {
		if distPath.Path == filePath {
			found = true
			dp = distPath
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("%w: unknown path: %s", ErrDetectionFailed, filePath)
	}

	for _, dist := range dp.Dists {
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
		{
			"/etc/os-release",
			`NAME="Oracle Linux Server"
VERSION="7.9"
ID="ol"
ID_LIKE="fedora"
VARIANT="Server"
VARIANT_ID="server"
VERSION_ID="7.9"
PRETTY_NAME="Oracle Linux Server 7.9"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:oracle:linux:7:9:server"
HOME_URL="https://linux.oracle.com/"
BUG_REPORT_URL="https://bugzilla.oracle.com/"

ORACLE_BUGZILLA_PRODUCT="Oracle Linux 7"
ORACLE_BUGZILLA_PRODUCT_VERSION=7.9
ORACLE_SUPPORT_PRODUCT="Oracle Linux"
ORACLE_SUPPORT_PRODUCT_VERSION=7.9`,
			info.Os{Family: runtime.GOOS, Arch: runtime.GOARCH, Distribution: "OracleLinux", Version: "7.9"},
			false,
		},
		{
			"/etc/os-release",
			`NAME="Red Hat Enterprise Linux"
VERSION="8.7 (Ootpa)"
ID="rhel"
ID_LIKE="fedora"
VERSION_ID="8.7"
PLATFORM_ID="platform:el8"
PRETTY_NAME="Red Hat Enterprise Linux 8.7 (Ootpa)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:redhat:enterprise_linux:8::baseos"
HOME_URL="https://www.redhat.com/"
DOCUMENTATION_URL="https://access.redhat.com/documentation/red_hat_enterprise_linux/8/"
BUG_REPORT_URL="https://bugzilla.redhat.com/"

REDHAT_BUGZILLA_PRODUCT="Red Hat Enterprise Linux 8"
REDHAT_BUGZILLA_PRODUCT_VERSION=8.7
REDHAT_SUPPORT_PRODUCT="Red Hat Enterprise Linux"
REDHAT_SUPPORT_PRODUCT_VERSION="8.7"`,
			info.Os{Family: runtime.GOOS, Arch: runtime.GOARCH, Distribution: "RedHat", Version: "8.7"},
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
			assert.Equal(t, tt.want, *got)
		})
	}
}
