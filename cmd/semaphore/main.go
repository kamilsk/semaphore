package main

import (
	"io"
	"os"
	"runtime"
	"text/template"
)

const (
	success = 0
	failure = 1
)

func main() { application{Args: os.Args, Stderr: os.Stderr, Stdout: os.Stdout, Shutdown: os.Exit}.Run() }

type application struct {
	Args           []string
	Stderr, Stdout io.Writer
	Shutdown       func(code int)
}

// Run executes the application logic.
func (app application) Run() {
	var command Command

	base := &BaseCommand{BinName: app.Args[0]}
	commands := Commands{
		&CreateCommand{BaseCommand: base,
			CmdName: "create", Capacity: runtime.GOMAXPROCS(0)},
		&AddCommand{BaseCommand: base,
			CmdName: "add"},
		&WaitCommand{BaseCommand: base,
			CmdName: "wait", Output: app.Stdout, Template: template.Must(template.New("report").Parse(DefaultReport))},
	}
	help := &HelpCommand{BaseCommand: base,
		CmdName: "help", Commit: commit, BuildDate: date, Version: version,
		Compiler: runtime.Compiler, Platform: runtime.GOOS + "/" + runtime.GOARCH, GoVersion: runtime.Version(),
		Commands: commands, Output: app.Stderr}
	commands = append(commands, help)

	if command, help.Error = commands.Parse(app.Args[1:]); help.Error != nil {
		if help.Do() != nil {
			app.Shutdown(failure)
			return
		}
		app.Shutdown(success)
		return
	}
	if help.Error = command.Do(); help.Error != nil {
		_ = help.Do()
		app.Shutdown(failure)
		return
	}
	app.Shutdown(success)
}
