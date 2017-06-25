package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	Version string
)

func main() {
	flag.Parse()
	fmt.Println(strings.Join(flag.Args(), ", "))
}
