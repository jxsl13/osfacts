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

func findSemanticVersion(content string) (string, error) {
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
		semVersionStr := semVersion.String()

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
		svStr := sv.String()

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