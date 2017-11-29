package main

import (
	"io"
	"os"
	"runtime"
	"text/template"
)

func main() { Main{Args: os.Args, Stderr: os.Stderr, Stdout: os.Stdout, Shutdown: os.Exit}.Exec() }

const (
	Success = 0
	Failed  = 1
)

// Main is a struct to use it in `main` function and tests.
type Main struct {
	Args           []string
	Stderr, Stdout io.Writer
	Shutdown       func(code int)
}

// Exec executes application logic.
func (m Main) Exec() {
	var command Command

	base := &BaseCommand{BinName: m.Args[0]}
	commands := Commands{
		&CreateCommand{BaseCommand: base.Copy(),
			CmdName: "create", Capacity: runtime.GOMAXPROCS(0)},
		&AddCommand{BaseCommand: base.Copy(),
			CmdName: "add"},
		&WaitCommand{BaseCommand: base.Copy(),
			CmdName: "wait", Output: m.Stdout, Template: template.Must(template.New("report").Parse(DefaultReport))},
	}
	help := &HelpCommand{BaseCommand: base.Copy(),
		CmdName: "help", Commit: commit, BuildDate: date, Version: version,
		Compiler: runtime.Compiler, Platform: runtime.GOOS + "/" + runtime.GOARCH, GoVersion: runtime.Version(),
		Commands: commands, Output: m.Stderr}
	commands = append(commands, help)

	if command, help.Error = commands.Parse(m.Args[1:]); help.Error != nil {
		if help.Do() != nil {
			m.Shutdown(Failed)
			return
		}
		m.Shutdown(Success)
		return
	}
	if help.Error = command.Do(); help.Error != nil {
		help.Do()
		m.Shutdown(Failed)
		return
	}
	m.Shutdown(Success)
	return
}
