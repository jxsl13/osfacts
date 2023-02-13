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

### Info

requires Linux kernel version 2.6.32

- SLES >= 11