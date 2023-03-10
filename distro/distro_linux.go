package distro

import (
	"fmt"
)

func detect() (*Info, error) {

	for _, distPath := range newPaths() {
		exists, err := existsWithSize(distPath.Path, false)
		if !exists || err != nil {
			continue
		}

		fileContent, err := getFileContent(distPath.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to open %s: %w", distPath.Path, err)
		}

		for _, dist := range distPath.Dists {
			osInfo, err := dist.Parse(distPath.Path, fileContent)
			if err != nil {
				continue
			}
			return osInfo, nil
		}
	}

	return nil, ErrDetectionFailed
}

// where could we use this distribution family
// https://github.com/ansible/ansible/blob/devel/lib/ansible/module_utils/facts/system/distribution.py#L512-L536

func newPaths() []distPath {

	defaultKeys := []string{
		"VERSION",
		"VERSION_ID",
		"PRETTY_NAME",
	}

	// order of these files is important
	// oracle linux > red hat
	return []distPath{
		{
			Path: "/etc/os-release",
			Dists: []distribution{
				{
					Name:      "Amazon",
					ParseFunc: parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:        "SUSE",
					SearchNames: []string{"SUSE", "suse"},
					ParseFunc:   parserFindEnvNameAndSemanticVersionKeys("NAME", "VERSION_ID"),
				},
				{
					Name:        "Debian",
					SearchNames: []string{"Debian", "Raspbian"},
					ParseFunc:   parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:        "Cumulus Linux",
					SearchNames: []string{"Cumulus"},
					ParseFunc:   parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:        "Linux Mint",
					SearchNames: []string{"Mint"},
					ParseFunc:   parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:      "CentOS",
					ParseFunc: parseCentOSDistFile,
				},
				{
					Name: "Uos", SearchNames: []string{"UOS", "Uos", "uos"},
					ParseFunc: parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name: "Deepin", SearchNames: []string{"Deepin", "deepin"},
					ParseFunc: parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:      "Ubuntu",
					ParseFunc: parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:      "Oracle Linux",
					ParseFunc: parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:      "Devuan",
					ParseFunc: parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:      "Arch Linux",
					ParseFunc: parserFindEnvVersionKey("BUILD_ID"),
				},
				{
					Name:        "RHEL",
					SearchNames: []string{"Red Hat"},
					ParseFunc:   parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					ParseFunc: parserFindEnvNameAndSemanticVersionKeys("NAME", defaultKeys...), // fallback
				},
			},
		},
		{
			Path: "/etc/altlinux-release",
			Dists: []distribution{
				{
					Name:        "ALT Linux",
					SearchNames: []string{"ALT"},
				}},
		},
		{
			Path: "/etc/oracle-release",
			Dists: []distribution{
				{
					Name: "Oracle Linux",
				},
			},
		},
		{
			Path: "/etc/slackware-version",
			Dists: []distribution{
				{
					Name:      "Slackware",
					ParseFunc: parserFindSemanticVersion,
				},
			},
		},
		{
			Path: "/etc/centos-release",
			Dists: []distribution{
				{
					Name:      "CentOS",
					ParseFunc: parserFindSemanticVersion,
				},
			},
		},
		{
			Path: "/etc/openwrt_release",
			Dists: []distribution{
				{
					Name:      "OpenWrt",
					ParseFunc: parserFindEnvSemanticVersionKeys("DISTRIB_RELEASE"),
				},
			},
		},
		{
			Path: "/etc/debian_version",
			Dists: []distribution{
				{
					Name:        "Debian",
					SearchNames: []string{"Debian", "Raspbian"},
					ParseFunc:   parserFindSemanticVersion,
				},
			},
		},
		{
			Path: "/etc/system-release",
			Dists: []distribution{
				{
					Name:      "Amazon",
					ParseFunc: parserFindSemanticVersion,
				},
			},
		},
		{
			Path: "/etc/alpine-release",
			Dists: []distribution{
				{
					Name:        "Alpine Linux",
					SearchNames: []string{"Alpine"},
					ParseFunc:   parserFindSemanticVersion,
				},
			},
		},
		{
			Path: "/etc/SuSE-release",
			Dists: []distribution{
				{
					Name:      "SUSE",
					ParseFunc: parseSuseReleaseDistFile,
				},
			},
		},
		{
			Path: "/etc/gentoo-release",
			Dists: []distribution{
				{
					Name: "Gentoo",
				},
			},
		},
		{
			Path: "/etc/lsb-release",
			Dists: []distribution{
				{
					Name:      "Debian",
					ParseFunc: parserFindEnvSemanticVersionKeys(defaultKeys...),
				},
				{
					Name:      "Mandriva",
					ParseFunc: parserFindEnvSemanticVersionKeys("DISTRIB_RELEASE"),
				},
			},
		},
		{
			Path: "/etc/sourcemage-release",
			Dists: []distribution{
				{
					Name: "SMGL", SearchNames: []string{"Source Mage GNU/Linux"},
				},
			},
		},
		{
			Path: "/usr/lib/os-release",
			Dists: []distribution{
				{
					Name:        "ClearLinux",
					SearchNames: []string{"Clear Linux"},
				},
			},
		},
		{
			Path: "/etc/redhat-release",
			Dists: []distribution{
				{
					Name:        "RHEL",
					SearchNames: []string{"Red Hat"},
					ParseFunc:   parserFindSemanticVersion,
				},
			},
		},
	}
}
