package main

import (
	"io"
	"os"
	"runtime"
	"text/template"
)

const (
	Success = 0
	Failed  = 1
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
		&CreateCommand{BaseCommand: base.Copy(),
			CmdName: "create", Capacity: runtime.GOMAXPROCS(0)},
		&AddCommand{BaseCommand: base.Copy(),
			CmdName: "add"},
		&WaitCommand{BaseCommand: base.Copy(),
			CmdName: "wait", Output: app.Stdout, Template: template.Must(template.New("report").Parse(DefaultReport))},
	}
	help := &HelpCommand{BaseCommand: base.Copy(),
		CmdName: "help", Commit: commit, BuildDate: date, Version: version,
		Compiler: runtime.Compiler, Platform: runtime.GOOS + "/" + runtime.GOARCH, GoVersion: runtime.Version(),
		Commands: commands, Output: app.Stderr}
	commands = append(commands, help)

	if command, help.Error = commands.Parse(app.Args[1:]); help.Error != nil {
		if help.Do() != nil {
			app.Shutdown(Failed)
			return
		}
		app.Shutdown(Success)
		return
	}
	if help.Error = command.Do(); help.Error != nil {
		help.Do()
		app.Shutdown(Failed)
		return
	}
	app.Shutdown(Success)
	return
}
