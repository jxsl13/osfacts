package distro

import (
	"fmt"
	"strings"
)

func mustContainOneOf(str string, substr ...string) (string, error) {
	if len(substr) == 0 {
		panic("substr must not be empty")
	}
	for _, s := range substr {
		if strings.Contains(str, s) {
			return s, nil
		}
	}

	return "", fmt.Errorf("string does not contain any of %v", substr)
}
