package common

import (
	"regexp"
	"strings"
)

var reSysctlSplit = regexp.MustCompile(`\s?=\s?|: `)

func GetSysctlCmd(args ...string) (*Command, error) {
	cmd, err := GetCommand("sysctl", args...)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func GetSysctl(args ...string) (map[string]string, error) {
	cmd, err := GetSysctlCmd(args...)
	if err != nil {
		return nil, err
	}
	lines, err := cmd.OutputLines()
	if err != nil {
		return nil, err
	}
	sysctl := make(map[string]string, len(lines))

	key := ""
	value := ""

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.HasPrefix(line, " ") {
			// handle multiline values, they will not have a starting key
			// Add the newline back in so people can split on it to parse
			// lines if they need to.
			value += "\n" + line
			continue
		}

		if key != "" {
			sysctl[key] = strings.TrimSpace(value)
		}

		parts := reSysctlSplit.Split(line, 2)
		switch len(parts) {
		case 0, 1:
			continue
		default:
			key = parts[0]
			value = parts[1]
		}
	}

	if key != "" {
		sysctl[key] = strings.TrimSpace(value)
	}

	return sysctl, nil
}
