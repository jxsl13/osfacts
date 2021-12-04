package common

import (
	"strings"
)

func GetBestParsableLocale(preferences ...string) (string, error) {
	cmd, err := GetCommand("locale", "-a")
	if err != nil {
		return "", err
	}

	if len(preferences) == 0 {
		preferences = []string{"C.utf8", "en_US.utf8", "C", "POSIX"}
	}

	available, err := cmd.OutputLines()
	if err != nil {
		return "", err
	}

	found := "C"
	if len(available) > 0 {
		for _, pref := range preferences {
			for _, avail := range available {
				if strings.Contains(avail, pref) {
					found = pref
					break
				}
			}
		}
	}
	return found, nil
}
