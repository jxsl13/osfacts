package distro

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/jxsl13/osfacts/internal"
)

var (
	macOSXConstraint *semver.Constraints
	osxContraint     *semver.Constraints
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

func detect() (*Info, error) {
	cmd, err := internal.NewCommand("/usr/bin/sw_vers", "-productVersion")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	output, err := cmd.OutputString()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}
	osInfo := newInfo()

	version, err := findSemVer(output)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDetectionFailed, err)
	}

	return osInfo.update(macOSName(version), version.String()), nil
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
