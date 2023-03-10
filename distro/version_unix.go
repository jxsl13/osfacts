//go:build !aix

package distro

import (
	"errors"
	"regexp"

	"github.com/Masterminds/semver"
)

var (
	ErrVersionNotFound = errors.New("no version found")
	semverRegex        = regexp.MustCompile(semver.SemVerRegex)
)

func findSemanticVersionString(content string) (string, error) {
	version := ""
	matches := semverRegex.FindAllString(content, -1)
	if len(matches) == 0 {
		return "", ErrVersionNotFound
	}

	for _, v := range matches {
		semVersion, err := semver.NewVersion(v)
		if err != nil {
			continue
		}
		semVersionStr := semVersion.Original()

		if len(semVersionStr) > len(version) {
			version = semVersionStr
		}
	}

	if version == "" {
		return "", ErrVersionNotFound
	}

	return version, nil
}

func findSemVer(content string) (*semver.Version, error) {
	var version *semver.Version
	vs := ""
	matches := semverRegex.FindAllString(content, -1)
	if len(matches) == 0 {
		return nil, ErrVersionNotFound
	}

	for _, v := range matches {
		sv, err := semver.NewVersion(v)
		if err != nil {
			continue
		}
		svStr := sv.Original()

		if len(svStr) > len(vs) {
			vs = svStr
			version = sv
		}
	}

	if vs == "" {
		return nil, ErrVersionNotFound
	}

	return version, nil
}
