package distro

import (
	"fmt"

	"github.com/jxsl13/osfacts/info"
)

// where could we use this distribution family
// https://github.com/ansible/ansible/blob/devel/lib/ansible/module_utils/facts/system/distribution.py#L512-L536

func newDistMap() map[string][]distribution {

	return map[string][]distribution{
		"/etc/altlinux-release": {{Name: "Altlinux", SearchNames: []string{"ALT"}}},
		"/etc/oracle-release":   {{Name: "OracleLinux", SearchNames: []string{"Oracle Linux"}}},
		"/etc/slackware-version": {
			{
				Name:      "Slackware",
				ParseFunc: parserFindSemanticVersion,
			},
		},
		"/etc/centos-release": {
			{
				Name:      "CentOS",
				ParseFunc: parseCentOSDistFile,
			},
		},
		"/etc/redhat-release": {
			{
				Name:        "RedHat",
				SearchNames: []string{"Red Hat"},
				ParseFunc:   parserFindSemanticVersion,
			},
		},
		"/etc/openwrt_release": {
			{
				Name:      "OpenWrt",
				ParseFunc: parserFindEnvSemanticVersionKeys("DISTRIB_RELEASE"),
			},
		},
		"/etc/debian_version": {
			{
				Name:        "Debian",
				SearchNames: []string{"Debian", "Raspbian"},
				ParseFunc:   parserFindSemanticVersion,
			},
		},
		"/etc/os-release": {
			{
				Name:      "Amazon",
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				Name:        "SUSE",
				SearchNames: []string{"SUSE", "suse"},
				ParseFunc:   parserFindEnvNameAndSemanticVersionKeys("NAME", "VERSION_ID"),
			},
			{
				Name:        "Debian",
				SearchNames: []string{"Debian", "Raspbian"},
				ParseFunc:   parserFindEnvSemanticVersionKeys(),
			},
			{
				Name: "Cumulus", Alias: "Cumulus Linux",
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				Name: "Mint", Alias: "Linux Mint",
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				Name: "Uos", SearchNames: []string{"UOS", "Uos", "uos"},
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				Name: "Deepin", SearchNames: []string{"Deepin", "deepin"},
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				Name:      "Ubuntu",
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				Name:      "Devuan",
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				Name:      "Archlinux",
				Alias:     "Arch Linux",
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				ParseFunc: parserFindEnvNameAndSemanticVersionKeys("NAME"), // fallback
			},
		},
		"/etc/system-release": {
			{
				Name:      "Amazon",
				ParseFunc: parserFindSemanticVersion,
			},
		},
		"/etc/alpine-release": {

			{
				Name:      "Alpine",
				ParseFunc: parserFindSemanticVersion,
			},
		},
		/*"/etc/arch-release": {
			{Name: "Archlinux", Alias: "Arch Linux"},
		},*/
		"/etc/SuSE-release": {
			{
				Name:      "SUSE",
				ParseFunc: parseSuseReleaseDistFile,
			},
		},
		"/etc/gentoo-release": {{Name: "Gentoo"}},
		"/etc/lsb-release": {
			{
				Name:      "Debian",
				ParseFunc: parserFindEnvSemanticVersionKeys(),
			},
			{
				Name:      "Mandriva",
				ParseFunc: parserFindEnvSemanticVersionKeys("DISTRIB_RELEASE"),
			},
		},
		"/etc/sourcemage-release": {
			{
				Name: "SMGL", SearchNames: []string{"Source Mage GNU/Linux"},
			},
		},
		"/usr/lib/os-release": {
			{Name: "ClearLinux", SearchNames: []string{"Clear Linux"}},
		},
	}

}

func detect() (*info.Os, error) {

	for filePath, dists := range newDistMap() {
		exists, err := existsWithSize(filePath, false)
		if !exists || err != nil {
			continue
		}

		fileContent, err := getFileContent(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", filePath, err)
		}

		for _, dist := range dists {
			osInfo, err := dist.Parse(filePath, fileContent)
			if err != nil {
				continue
			}
			return osInfo, nil
		}
	}

	return nil, ErrDetectionFailed
}
