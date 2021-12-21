package common

import "strings"

func InSlice(s string, slice []string) bool {
	for _, val := range slice {
		if val == s {
			return true
		}
	}
	return false
}

func OneOf(s string, values ...string) bool {
	for _, v := range values {
		if s == v {
			return true
		}
	}
	return false
}

func StartsWithOneOf(s string, values ...string) bool {
	for _, v := range values {
		if strings.HasPrefix(s, v) {
			return true
		}
	}
	return false
}
