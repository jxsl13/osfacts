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
