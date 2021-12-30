package common

import (
	"errors"
	"strings"
)

var (
	ErrTooFewParts = errors.New("too few parts")
)

func SplitLines(text string) []string {
	lines := strings.Split(text, "\n")
	return lines
}

func SplitLineIntoAtLeast(line, sep string, parts int) ([]string, error) {
	s := strings.SplitN(line, sep, parts)
	if len(s) < parts {
		return nil, ErrTooFewParts
	}

	return s, nil
}
