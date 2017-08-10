package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	flag.Parse()
	fmt.Println(strings.Join(flag.Args(), ", "))
	fmt.Println(commit, date, version)
}
