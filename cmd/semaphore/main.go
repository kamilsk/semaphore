package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	bc := &BaseCommand{BinName: os.Args[0]}
	commands := Commands{
		&CreateCommand{BaseCommand: bc,
			CmdName: "create", Capacity: runtime.GOMAXPROCS(0)},
		&AddCommand{BaseCommand: bc,
			CmdName: "add"},
		&WaitCommand{BaseCommand: bc,
			CmdName: "wait", Stdout: os.Stdout, Stderr: os.Stderr},
	}

	command, err := commands.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := command.Do(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
