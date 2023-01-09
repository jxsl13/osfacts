package distro

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jxsl13/osfacts/info"
)

var (
	osDistList = []distribution{
		{Path: "/etc/os-release", Name: "Amazon"},
		{Path: "/etc/os-release", Name: "Archlinux"},
		{Path: "/etc/os-release", Name: "Debian"},
		{Path: "/etc/os-release", Name: "Flatcar"},
		{Path: "/etc/os-release", Name: "SUSE"},
		{Path: "/etc/altlinux-release", Name: "Altlinux"},
		{Path: "/etc/oracle-release", Name: "OracleLinux"},
		{Path: "/etc/slackware-version", Name: "Slackware"},
		{Path: "/etc/centos-release", Name: "CentOS"},
		{Path: "/etc/redhat-release", Name: "RedHat"},
		{Path: "/etc/openwrt_release", Name: "OpenWrt"},
		{Path: "/etc/system-release", Name: "Amazon"},
		{Path: "/etc/alpine-release", Name: "Alpine"},
		{Path: "/etc/SuSE-release", Name: "SUSE"},
		{Path: "/etc/gentoo-release", Name: "Gentoo"},
		{Path: "/etc/lsb-release", Name: "Debian"},
		{Path: "/etc/lsb-release", Name: "Mandriva"},
		{Path: "/etc/sourcemage-release", Name: "SMGL"},
		{Path: "/usr/lib/os-release", Name: "ClearLinux"},
		{Path: "/etc/os-release", Name: "NA"},
	}

	searchStringMap = map[string]string{
		"OracleLinux": "Oracle Linux",
		"RedHat":      "Red Hat",
		"Altlinux":    "ALT",
		"SMGL":        "Source Mage GNU/Linux",
	}

	releaseAlias = map[string]string{
		"Archlinux": "Arch Linux",
	}

	parsers = map[string]fileParseFunc{
		"Amazon":    ParseAmazonDistFile,
		"Alpine":    ParseAlpineDistFile,
		"Debian":    ParseDebianDistFile,
		"OpenWrt":   ParseOpenWrtDistFile,
		"Slackware": ParseSlackwareDistFile,
		"SUSE":      ParseSUSEDistFile,
		"NA":        ParseFallbackDistFile,
	}
)

func detect() (*info.Os, error) {
	for _, dist := range osDistList {
		content, err := dist.Content()
		if err != nil {
			continue
		}
		info, err := parseDistFile(dist, content)
		if err != nil {
			continue
		} else {
			return info, nil
		}
	}

	return nil, ErrDetectionFailed
}

func parseDistFile(dist distribution, fileContent string) (*info.Os, error) {
	osInfo := info.NewOs()

	if searchString, found := searchStringMap[dist.Name]; found {
		if strings.Contains(fileContent, searchString) {
			osInfo.Distribution = dist.Name
		} else {
			tokens := strings.SplitN(fileContent, " ", 2)
			if len(tokens) != 2 {
				return nil, fmt.Errorf("invalid distribution release file content (%s): %s", dist.Name, dist.Path)
			}
			// use first string from /etc/xyz file
			osInfo.Distribution = strings.TrimSpace(stripQuotes(tokens[0]))
		}

		if alias, found := releaseAlias[dist.Name]; found && strings.Contains(fileContent, alias) {
			osInfo.Distribution = alias
		}

		versionString, err := findSemanticVersion(fileContent)
		if err != nil {
			return nil, err
		}
		osInfo.Version = versionString
		return osInfo, nil
	}

	parser, found := parsers[dist.Name]
	if !found {
		return nil, errors.New("distribution could not be detected")
	}

	return parser(dist, fileContent)
}
