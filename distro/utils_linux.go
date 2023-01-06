package distro

import (
	"errors"
	"os"
	"strings"
)

var (
	quoteReplacer = strings.NewReplacer(
		`"`, "",
		`'`, "",
		`\\`, "",
	)
	ErrKeyNotFound = errors.New("key not found")
)

func stripQuotes(src string) string {
	return quoteReplacer.Replace(src)
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
