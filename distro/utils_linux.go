package distro

import (
	"errors"
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
