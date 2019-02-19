package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

var (
	errNotProvided = fmt.Errorf("command not provided")
	errNotFound    = fmt.Errorf("command not found")
	errExecution   = fmt.Errorf("command execution error")
)

// Command defines behavior to interact with user input.
type Command interface {
	// FlagSet should return configured FlagSet to handle command arguments.
	FlagSet() *flag.FlagSet
	// Name should return a command name.
	Name() string
	// Desc should return a command description.
	Desc() string
	// Do should exec a command.
	Do() error
}

// Commands is a container provides the method to search an appropriate command.
type Commands []Command

// Parse parses the arguments and searches an appropriate command for them.
func (l Commands) Parse(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, errNotProvided
	}
	cmdName := args[0]
	if _, found := map[string]struct{}{"-h": {}, "-help": {}, "--help": {}}[cmdName]; found {
		return nil, flag.ErrHelp
	}
	for _, cmd := range l {
		if cmd.Name() == cmdName {
			return cmd, errors.Wrapf(cmd.FlagSet().Parse(args[1:]), "invalid arguments for command %s", cmd.Name())
		}
	}
	return nil, errNotFound
}

// BaseCommand contains general fields for other commands.
type BaseCommand struct {
	Debug    bool
	BinName  string
	FileName string
	Mode     flag.ErrorHandling
}

// FlagSet creates and configures new general FlagSet.
func (c *BaseCommand) FlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, c.Mode)
	fs.BoolVar(&c.Debug, "debug", false, "show error stack trace")
	fs.StringVar(&c.FileName, "filename", filepath.Join(os.TempDir(), c.BinName+".json"),
		"an absolute path to semaphore context")
	return fs
}

// CreateCommand is a command to init a semaphore context.
type CreateCommand struct {
	*BaseCommand
	CmdName  string
	Capacity int
	Flags    *flag.FlagSet
}

// FlagSet returns a configured FlagSet to handle CreateCommand arguments.
func (c *CreateCommand) FlagSet() *flag.FlagSet {
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
	}
	return c.Flags
}

// Name returns a CreateCommand name.
func (c *CreateCommand) Name() string {
	return c.CmdName
}

// Desc returns a CreateCommand description.
func (c *CreateCommand) Desc() string {
	return "is a command to init a semaphore context"
}

// Do creates a file to store a semaphore context.
func (c *CreateCommand) Do() error {
	var err error

	args := c.FlagSet().Args()
	capacity := c.Capacity
	if len(args) > 0 {
		if capacity, err = strconv.Atoi(args[0]); err != nil || capacity < 1 {
			return errors.Wrapf(err, "invalid capacity: capacity must be a valid integer greater than zero")
		}
	}

	file, err := os.Create(c.FileName)
	if err != nil {
		return errors.Wrapf(err, "could not create a file %s", c.FileName)
	}

	task := Task{Capacity: capacity}
	if err := json.NewEncoder(file).Encode(task); err != nil {
		return errors.Wrapf(err, "could not store a context %+v into a file %s", task, c.FileName)
	}

	return nil
}

// AddCommand is a command to add a job into a semaphore context.
type AddCommand struct {
	*BaseCommand
	CmdName string
	Edit    bool
	Command []string
	Flags   *flag.FlagSet
}

// FlagSet returns configured FlagSet to handle AddCommand arguments.
func (c *AddCommand) FlagSet() *flag.FlagSet {
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
		c.Flags.BoolVar(&c.Edit, "edit", false, "switch to edit mode to read arguments from input (not implemented yet)")
	}
	return c.Flags
}

// Name returns an AddCommand name.
func (c *AddCommand) Name() string {
	return c.CmdName
}

// Desc returns an AddCommand description.
func (c *AddCommand) Desc() string {
	return "is a command to add a job into a semaphore context"
}

// Do adds a job into a semaphore context and stores it.
func (c *AddCommand) Do() error {
	if c.Edit {
		// TODO each new line from os.Stdin should be converted to Task
		_, _ = color.New(color.FgYellow).Fprintln(os.Stderr, "edit component is not ready yet")
	}

	args := c.FlagSet().Args()
	if len(args) == 0 {
		return fmt.Errorf("the add command requires arguments to create a job based on them")
	}

	file, err := os.OpenFile(c.FileName, os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrapf(err, "could not open a file %s. did you create it before?", c.FileName)
	}

	var task Task
	if err := json.NewDecoder(file).Decode(&task); err != nil {
		return errors.Wrapf(err, "could not restore a context from a file %s. is it a valid JSON?", c.FileName)
	}

	task.AddJob(Job{Name: args[0], Args: args[1:]})
	data, err := json.Marshal(task)
	if err != nil {
		return errors.Wrapf(err, "could not encode a context %+v into a JSON", task)
	}

	if _, err := file.WriteAt(data, 0); err != nil {
		return errors.Wrapf(err, "could not store a context %+v into a file %s", task, c.FileName)
	}

	return nil
}

// DefaultReport is a default template for report.
var DefaultReport = `
command: {{ .Name }} {{ range .Args }}{{ . }} {{ end }}
  error: {{ .Error }}
details: started at {{ .Start }}, finished at {{ .End }}, elapsed {{ .Elapsed }}
 stdout:

{{ .Stdout }}

 stderr:

{{ .Stderr }}
---
`

// ColoredOutput wraps another output and colorizes input data before to pass it.
type ColoredOutput struct {
	clr *color.Color
	dst io.Writer
}

// Write implements `io.Writer` interface.
func (o *ColoredOutput) Write(p []byte) (int, error) {
	return o.clr.Fprint(o.dst, string(p))
}

// LimitedOutput wraps another output and writes to it with specified velocity.
type LimitedOutput struct {
	dst   io.Writer
	speed int
}

// For sets limiter for passed io.Writer.
func (o *LimitedOutput) For(dst io.Writer) *LimitedOutput {
	o.dst = dst
	return o
}

// Write implements `io.Writer` interface.
func (o *LimitedOutput) Write(p []byte) (int, error) {
	if o.speed != 0 {
		pause := time.Second / time.Duration(o.speed)
		var (
			total, n int
			err      error
		)
		for _, r := range string(p) {
			if r == 0 {
				continue
			}
			n, err = o.dst.Write([]byte(string(r)))
			total += n
			if err != nil {
				break
			}
			time.Sleep(pause)
		}
		return total, err
	}
	return o.dst.Write(p)
}

// WaitCommand is a command to execute a semaphore task.
type WaitCommand struct {
	*BaseCommand
	CmdName  string
	Notify   bool
	Output   io.Writer
	Speed    int
	Template *template.Template
	Timeout  time.Duration
	Flags    *flag.FlagSet
}

// FlagSet returns a configured FlagSet to handle WaitCommand arguments.
func (c *WaitCommand) FlagSet() *flag.FlagSet {
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
		c.Flags.BoolVar(&c.Notify, "notify", false, "show notification at the end (not implemented yet)")
		c.Flags.DurationVar(&c.Timeout, "timeout", time.Minute, "timeout for task execution")
		c.Flags.IntVar(&c.Speed, "speed", 0, "a velocity of report output (characters per second)")
	}
	return c.Flags
}

// Name returns a WaitCommand name.
func (c *WaitCommand) Name() string {
	return c.CmdName
}

// Desc returns a WaitCommand description.
func (c *WaitCommand) Desc() string {
	return "is a command to execute a semaphore task"
}

// Do executes a semaphore task.
func (c *WaitCommand) Do() error {
	file, err := os.OpenFile(c.FileName, os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrapf(err, "could not open a file %s. did you create it before?", c.FileName)
	}

	var task Task
	if err := json.NewDecoder(file).Decode(&task); err != nil {
		return errors.Wrapf(err, "could not restore a context from a file %s. is it a valid JSON?", c.FileName)
	}
	if c.Timeout > 0 {
		task.Timeout = c.Timeout
	}

	var (
		bar              = pb.New(len(task.Jobs))
		results          = &Results{}
		red              = &ColoredOutput{clr: color.New(color.FgHiRed), dst: c.Output}
		limiter          = &LimitedOutput{speed: c.Speed}
		success, failure = 0, 0
		output           io.Writer
		start, end       time.Time
	)
	bar.Output = c.Output
	bar.ShowTimeLeft = false
	bar.Start()
	start = time.Now()
	for result := range task.Run() {
		success++
		if result.Error != nil {
			bar.Output = red
			failure++
			success--
		}
		bar.Increment()
		results.Append(result)
	}
	end = time.Now()
	bar.Finish()

	for _, result := range results.Sort() {
		output = c.Output
		if result.Error != nil {
			output = red
		}
		stdout, _ := ioutil.ReadAll(result.Stdout)
		stderr, _ := ioutil.ReadAll(result.Stderr)
		err = errors.Wrap(c.Template.Execute(limiter.For(output), struct {
			Name       string
			Args       []string
			Error      error
			Start, End string
			Elapsed    time.Duration
			Stdout     string
			Stderr     string
		}{
			Name:    result.Job.Name,
			Args:    result.Job.Args,
			Error:   result.Error,
			Start:   result.Start.Format("2006-01-02 15:04:05.99"),
			End:     result.End.Format("2006-01-02 15:04:05.99"),
			Elapsed: result.End.Sub(result.Start),
			Stdout:  string(stdout),
			Stderr:  string(stderr),
		}), "template execution")
	}

	output = c.Output
	if failure > 0 {
		output = red
	}
	_, _ = fmt.Fprintf(output, "total: %d; successful: %d; failed: %d; elapsed: %s \n",
		results.Len(), success, failure, end.Sub(start))

	if c.Notify {
		// TODO try to find or implement by myself
		// - https://github.com/variadico/noti
		// - https://github.com/jolicode/JoliNotif
		_, _ = color.New(color.FgYellow).Fprintln(os.Stderr, "notify component is not ready yet")
	}

	if failure > 0 {
		return errExecution
	}
	return err
}

// HelpCommand is command to show help message.
type HelpCommand struct {
	*BaseCommand
	CmdName                       string
	Commit, BuildDate, Version    string
	Compiler, Platform, GoVersion string
	Commands                      Commands
	Error                         error
	Output                        io.Writer
	Flags                         *flag.FlagSet
}

// FlagSet returns a configured FlagSet to handle HelpCommand arguments.
func (c *HelpCommand) FlagSet() *flag.FlagSet {
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
	}
	return c.Flags
}

// Name returns a HelpCommand name.
func (c *HelpCommand) Name() string {
	return c.CmdName
}

// Desc returns a HelpCommand description.
func (c *HelpCommand) Desc() string {
	return "is command to show help message"
}

// Do handles inner error and shows a specific message.
func (c *HelpCommand) Do() error {
	switch c.Error {
	case nil, errNotProvided, flag.ErrHelp:
		c.Usage()
		return nil
	case errNotFound:
		c.Usage()
		fallthrough
	case errExecution:
		return c.Error
	default:
		_, _ = color.New(color.FgRed).Fprintf(c.Output, "an error occurred: %v\n", c.Error)
		return c.Error
	}
}

// Usage shows help message.
func (c *HelpCommand) Usage() {
	_, _ = fmt.Fprintf(c.Output, `
Usage: %s COMMAND

Semaphore provides functionality to execute terminal commands in parallel.

`, c.BinName)

	if len(c.Commands) > 0 {
		_, _ = fmt.Fprintln(c.Output, "Commands:")
		for _, cmd := range c.Commands {
			_, _ = fmt.Fprintf(c.Output, "\n%s\t%s\n", cmd.Name(), cmd.Desc())
			fs := cmd.FlagSet()
			fs.SetOutput(c.Output)
			fs.PrintDefaults()
			_, _ = fmt.Fprintln(c.Output)
		}
	}

	_, _ = fmt.Fprintf(c.Output, "Version %s (commit: %s, build date: %s, go version: %s, compiler: %s, platform: %s)\n",
		c.Version, c.Commit, c.BuildDate, c.GoVersion, c.Compiler, c.Platform)
}
