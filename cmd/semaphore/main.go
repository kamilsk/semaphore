package main

import (
	"os"
	"runtime"
	"text/template"
)

func main() {
	var command Command

	base := &BaseCommand{BinName: os.Args[0]}
	commands := Commands{
		&CreateCommand{BaseCommand: base.Copy(),
			CmdName: "create", Capacity: runtime.GOMAXPROCS(0)},
		&AddCommand{BaseCommand: base.Copy(),
			CmdName: "add"},
		&WaitCommand{BaseCommand: base.Copy(),
			CmdName: "wait",
			Stdout:  os.Stdout, Stderr: os.Stderr,
			Template: template.Must(template.New("report").Parse(DefaultReport))},
	}
	help := &HelpCommand{BaseCommand: base.Copy(),
		CmdName: "help", Commit: commit, Date: date, Version: version, Commands: commands, Output: os.Stderr}
	commands = append(commands, help)

	if command, help.Error = commands.Parse(os.Args[1:]); help.Error != nil {
		if help.Do() != nil {
			os.Exit(1)
		}
		return
	}
	if help.Error = command.Do(); help.Error != nil {
		help.Do()
		os.Exit(1)
	}
}
