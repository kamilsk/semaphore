package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
)

/*
Prototype:
$ semaphore create 4
$ semaphore add -- docker build ...
$ semaphore add -- docker build ...
...
$ semaphore wait | semaphore wait --notify --timeout 1h
... show progress (colored output)
[==>........] 2/10

command `docker build ...`
output:
 ...

command...
*/
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	select {
	case <-c:
		cancel()
	case <-ctx.Done():
	}

	flag.Parse()
	fmt.Println(strings.Join(flag.Args(), ", "))
	fmt.Println(commit, date, version, os.TempDir())
}
