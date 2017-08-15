package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	base := &BaseCommand{BinName: os.Args[0]}
	commands := Commands{
		&CreateCommand{BaseCommand: base.Copy(),
			CmdName: "create", Capacity: runtime.GOMAXPROCS(0)},
		&AddCommand{BaseCommand: base.Copy(),
			CmdName: "add"},
		&WaitCommand{BaseCommand: base.Copy(),
			CmdName: "wait", Stdout: os.Stdout, Stderr: os.Stderr},
	}
	help := &HelpCommand{BaseCommand: base.Copy(),
		CmdName: "help", Commit: commit, Date: date, Version: version, Commands: []Command(commands), Output: os.Stderr}

	command, err := commands.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		help.Do()
		os.Exit(1)
	}
	if err := command.Do(); err != nil {
		fmt.Println(err)
		help.Do()
		os.Exit(1)
	}
}
