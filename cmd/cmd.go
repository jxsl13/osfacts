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
