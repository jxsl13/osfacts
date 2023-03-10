# osfacts

Small library that allows to detect system information like operating system, its version and distribution

### import

```shell
go get github.com/jxsl13/osfacts@latest
```

### example

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/jxsl13/osfacts/distro"
)

func main() {
	o, err := distro.Detect()
	if err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(o, "", " ")
	fmt.Println(string(data))
}
```

### requirements

requires Linux kernel version 2.6.32 or newer, which is the smallest compilation target for Go.
This prevents us from running on pretty old operating systems like SLES 10 and older.

- SLES >= 11


## `distro` package

This package supports detecting different os families, architectures, distributions and distribution versions.

### supportes OS families
The `distro` package supports these families and detecting their distro (or os) version.

- aix (does not have distributions)
- darwin (different os names depending on version)
- linux (see below)
- solaris (see below)
- windows 

### supported `linux` distributions

- Alpine Linux
- ALT Linux
- Amazon
- Arch Linux
- CentOS (and CentOS Stream, TencentOS)
- Clear Linux
- Cumulus Linux
- Debian
- Deepin
- Devuan
- Gentoo
- Linux Mint
- Mandriva
- OpenWrt
- Oracle Linux
- RHEL
- Slackware
- SMGL
- SUSE
- Ubuntu
- Uos

### supported `solaris` distributions

- NexentaOS
- OmniOS
- OpenIndiana
- SmartOS
- Solaris (= Oracle Solaris)

### `windows`

The windows versioning is the same as the version string provided by the `WinVer.exe`.

