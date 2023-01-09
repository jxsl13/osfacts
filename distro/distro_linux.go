package distro

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jxsl13/osfacts/info"
)

var (
	// TODO: distribution should be a public interface with a Detect method
	osDistList = []distribution{
		{Path: "/etc/altlinux-release", Name: "Altlinux"},
		{Path: "/etc/oracle-release", Name: "OracleLinux"},
		{Path: "/etc/slackware-version", Name: "Slackware"},
		{Path: "/etc/centos-release", Name: "CentOS"},
		{Path: "/etc/redhat-release", Name: "RedHat"},
		//{Path: "/etc/vmware-release", Name: "VMwareESX"}, // TODO: add custom detection
		{Path: "/etc/openwrt_release", Name: "OpenWrt"},
		{Path: "/etc/os-release", Name: "Amazon"},
		{Path: "/etc/system-release", Name: "Amazon"},
		{Path: "/etc/alpine-release", Name: "Alpine"},
		{Path: "/etc/arch-release", Name: "Archlinux"},
		// {Path: "/etc/os-release", Name: "Archlinux"}, // TODO: add custom detection
		{Path: "/etc/os-release", Name: "SUSE"},
		{Path: "/etc/SuSE-release", Name: "SUSE"},
		{Path: "/etc/gentoo-release", Name: "Gentoo"},
		{Path: "/etc/os-release", Name: "Debian"},
		{Path: "/etc/lsb-release", Name: "Debian"},
		{Path: "/etc/lsb-release", Name: "Mandriva"},
		{Path: "/etc/sourcemage-release", Name: "SMGL"},
		{Path: "/usr/lib/os-release", Name: "ClearLinux"},
		//{Path: "/etc/coreos/update.conf", Name: "Coreos"}, // TODO: need test environment
		// {Path: "/etc/os-release", Name: "Flatcar"},  // TODO: need test environment
		{Path: "/etc/os-release", Name: ""}, // fallback search
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
		"Amazon":     parseAmazonDistFile,
		"Alpine":     parseAlpineDistFile,
		"CentOS":     parseCentOSDistFile,
		"Debian":     parseDebianDistFile,
		"OpenWrt":    parseOpenWrtDistFile,
		"Slackware":  parseSlackwareDistFile,
		"SUSE":       parseSUSEDistFile,
		"Mandriva":   parseMandrivaDistFile,
		"ClearLinux": parseClearLinuxDistFile,
		"":           parseFallbackDistFile,
	}
)

func detect() (*info.Os, error) {
	var (
		osInfo = info.NewOs()
	)

	for _, dist := range osDistList {
		content, err := dist.Content()
		if err != nil {
			// file not found
			continue
		}
		err = parseDistFile(dist, content, osInfo)
		if err == nil {
			return osInfo, nil
		}
	}

	return nil, ErrDetectionFailed
}

func parseDistFile(dist distribution, fileContent string, osInfo *info.Os) error {

	if searchString, found := searchStringMap[dist.Name]; found {
		if strings.Contains(fileContent, searchString) {
			osInfo.Distribution = dist.Name
		} else {
			tokens := strings.SplitN(fileContent, " ", 2)
			if len(tokens) != 2 {
				return fmt.Errorf("invalid distribution release file content (%s): %s", dist.Name, dist.Path)
			}
			// use first string from /etc/xyz file
			osInfo.Distribution = strings.TrimSpace(stripQuotes(tokens[0]))
		}

		if alias, found := releaseAlias[dist.Name]; found && strings.Contains(fileContent, alias) {
			osInfo.Distribution = alias
		}

		versionString, err := findSemanticVersion(fileContent)
		if err != nil {
			return err
		}
		osInfo.Version = versionString
		return nil
	}

	parser, found := parsers[dist.Name]
	if !found {
		return errors.New("distribution could not be detected")
	}

	return parser(dist, fileContent, osInfo)
}
