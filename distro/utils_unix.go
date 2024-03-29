//go:build !aix

package distro

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/joho/godotenv"
)

var (
	ErrKeyNotFound       = errors.New("key not found")
	ErrInvalidFileFormat = errors.New("invalid file format")
)

func existsWithSize(filePath string, allowEmpty bool) (bool, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	if allowEmpty {
		return true, nil
	}

	if fi.Size() == 0 {
		return false, nil
	}

	return true, nil
}

func findEnvSemanticVersionInMap(envMap map[string]string, keys ...string) (string, error) {
	var (
		version  *semver.Version
		original string = ""
	)
	if len(keys) == 0 {
		// look through whole map
		for _, value := range envMap {
			v, err := findSemVer(value)
			if err != nil {
				continue
			}

			vOriginal := v.Original()

			if version == nil {
				version = v
				original = vOriginal
				continue
			}

			if len(vOriginal) > len(original) {
				version = v
				original = vOriginal
			} else if len(vOriginal) == len(original) && version.LessThan(v) {
				// 12.1.0 vs 12.1.2 (equal length but we want the higher version)
				version = v
				original = vOriginal
			}
		}
	} else {
		// look through keys only
		for _, key := range keys {
			value, found := envMap[key]
			if !found {
				continue
			}

			v, err := findSemVer(value)
			if err != nil {
				continue
			}

			vOriginal := v.Original()
			if version == nil {
				version = v
				original = vOriginal
				continue
			}

			if len(vOriginal) > len(original) {
				version = v
				original = vOriginal
			} else if len(vOriginal) == len(original) && version.LessThan(v) {
				version = v
				original = vOriginal
			}
		}
	}

	if version == nil {
		return "", ErrVersionNotFound
	}

	return original, nil
}

// set keys to look at for versions or keep empty to look at all keys
func findEnvSemanticVersion(content string, keys ...string) (string, error) {
	envMap, err := godotenv.Parse(strings.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidFileFormat, err)
	}
	return findEnvSemanticVersionInMap(envMap, keys...)
}

// reference: https://github.com/ansible/ansible/blob/616ad883addfc2ca776d8050dd2aa19ae6d1f1d0/lib/ansible/module_utils/facts/system/distribution.py

func getKey(m map[string]string, key string) (string, error) {
	value, found := m[key]
	if !found {
		return "", fmt.Errorf("%w: %q", ErrKeyNotFound, key)
	}
	return value, nil
}

func getEnvMap(content string) (map[string]string, error) {
	envMap, err := godotenv.Parse(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidFileFormat, err)
	}
	return envMap, nil
}

func getFileContent(path string) (string, error) {
	found, err := existsWithSize(path, false)
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	if !found {
		return "", fmt.Errorf("%s: not found or empty", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	return string(data), nil
}
