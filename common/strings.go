package common

import (
	"errors"
	"strings"
)

var (
	ErrTooFewParts = errors.New("too few parts")
)

func SplitLines(text string) []string {
	return strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
}

func SplitLineIntoAtLeast(line, sep string, parts int) ([]string, error) {
	s := strings.SplitN(line, sep, parts)
	if len(s) < parts {
		return nil, ErrTooFewParts
	}

	return s, nil
}
