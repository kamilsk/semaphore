package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
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
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
	}()

	commands := Commands{
		&CreateCommand{BaseCommand: BaseCommand{ID: "create"}, Filename: filepath.Join(os.TempDir(), os.Args[0]+".json")},
		&AddCommand{BaseCommand: BaseCommand{ID: "add"}},
		&WaitCommand{BaseCommand: BaseCommand{ID: "wait"}},
	}
	command, err := commands.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
	command.Do()

	fmt.Println(commit, date, version)
}
