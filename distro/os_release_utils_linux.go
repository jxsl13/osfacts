package distro

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ErrInvalidFileFormat = errors.New("invalid file format")
)

func findOsReleaseSemanticVersionInMap(envMap map[string]string, keys ...string) (string, error) {
	version := ""
	if len(keys) == 0 {
		// look through whole map
		for _, value := range envMap {
			v, err := findSemanticVersion(value)
			if err != nil {
				continue
			}

			if len(v) > len(version) {
				version = v
			}
		}
	} else {
		// look through keys only
		for _, key := range keys {
			value, found := envMap[key]
			if !found {
				continue
			}
			v, err := findSemanticVersion(value)
			if err != nil {
				continue
			}

			if len(v) > len(version) {
				version = v
			}
		}
	}

	if version == "" {
		return "", ErrVersionNotFound
	}

	return version, nil
}

// set keys to look at for versions or keep empty to look at all keys
func findOsReleaseSemanticVersion(content string, keys ...string) (string, error) {
	envMap, err := godotenv.Parse(strings.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidFileFormat, err)
	}
	return findOsReleaseSemanticVersionInMap(envMap, keys...)
}

// reference: https://github.com/ansible/ansible/blob/616ad883addfc2ca776d8050dd2aa19ae6d1f1d0/lib/ansible/module_utils/facts/system/distribution.py

func getKey(m map[string]string, key string) (string, error) {
	value, found := m[key]
	if !found {
		return "", fmt.Errorf("%w: %q", ErrKeyNotFound, key)
	}
	return value, nil
}

func getOsReleaseMap(content string) (map[string]string, error) {
	envMap, err := godotenv.Parse(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidFileFormat, err)
	}
	return envMap, nil
}
