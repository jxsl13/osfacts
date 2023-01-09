package distro

import (
	"fmt"

	"github.com/jxsl13/osfacts/common"
	"github.com/jxsl13/osfacts/info"
)

var (
	macOSXConstraint *semver.Constraint
	osxContraint *semver.Constraint
	macOSConstraint *semver.Constraint
)

func init() {
	var err error
	macOSXConstraint, err = semver.NewConstraint("< 10.8")
    if err != nil {
        panic(err)
    }

	osxContraint, err = semver.NewConstraint("< 10.12")
	if err != nil {
        panic(err)
    }
}

// reference: https://n8felton.wordpress.com/2022/01/28/macos-version-naming-conventions/
func macOSName(version *semver.Version) string {
	if macOSXConstraint.Check(version) {
		// < 10.8
		return "Mac OS X"
	} else if osxContraint.Check(version) {
		// < 10.12
		return "OS X"
	} else {
		// >= 10.12
		return "macOS"
	}

}

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

	version, err := findSemVer(output)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}

	osInfo.Distribution = macOSName(version)
	osInfo.Version = version.String()


	if

	return osInfo, nil
}
