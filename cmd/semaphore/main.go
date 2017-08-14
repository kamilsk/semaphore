package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
)

func main() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
	}()

	filename := filepath.Join(os.TempDir(), os.Args[0]+".json")
	commands := Commands{
		&CreateCommand{
			BaseCommand: BaseCommand{ID: "create", Filename: filename},
			Capacity:    runtime.GOMAXPROCS(0)},
		&AddCommand{
			BaseCommand: BaseCommand{ID: "add", Filename: filename}},
		&WaitCommand{
			BaseCommand: BaseCommand{Bin: os.Args[0], ID: "wait", Filename: filename},
			Stdout:      os.Stdout, Stderr: os.Stderr},
	}

	command, err := commands.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
	command.Do()
}
