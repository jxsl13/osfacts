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
	quoteReplacer = strings.NewReplacer(
		`"`, "",
		`'`, "",
		`\\`, "",
	)
	ErrKeyNotFound       = errors.New("key not found")
	ErrInvalidFileFormat = errors.New("invalid file format")
)

func stripQuotes(src string) string {
	return quoteReplacer.Replace(src)
}

func unique[T comparable](values []T) []T {
	m := make(map[T]struct{}, len(values))
	for _, v := range values {
		m[v] = struct{}{}
	}

	result := make([]T, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

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

			if version == nil {
				version = v
				original = v.String()
				continue
			}

			vOriginal := v.Original()

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

			if version == nil {
				version = v
				original = v.String()
				continue
			}

			vOriginal := v.Original()

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
