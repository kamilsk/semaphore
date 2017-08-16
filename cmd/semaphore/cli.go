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
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	NotProvided = fmt.Errorf("command not provided")
	NotFound    = fmt.Errorf("command not found")
)

// Command defines behavior to interact with user input.
type Command interface {
	FlagSet() *flag.FlagSet
	Name() string
	Desc() string
	Do() error
}

// Commands is a container which provides the method to search an appropriate command.
type Commands []Command

// Parse parses the arguments and searches an appropriate command for them.
func (l Commands) Parse(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, NotProvided
	}
	cmdName := args[0]
	for _, cmd := range l {
		if cmd.Name() == cmdName {
			return cmd, errors.WithMessage(cmd.FlagSet().Parse(args[1:]),
				fmt.Sprintf("invalid arguments for command %s", cmd.Name()))
		}
	}
	return nil, NotFound
}

// BaseCommand contains general properties for other commands.
type BaseCommand struct {
	BinName  string
	FileName string
	Mode     flag.ErrorHandling
	Flags    *flag.FlagSet
}

// Copy returns a copy of current BaseCommand.
func (c *BaseCommand) Copy() *BaseCommand {
	n := *c
	return &n
}

// FlagSet creates and configures new general FlagSet.
func (c *BaseCommand) FlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, c.Mode)
	value := filepath.Join(os.TempDir(), c.BinName+".json")
	fs.StringVar(&c.FileName, "filename", value, "an absolute path to semaphore's context")
	return fs
}

// ~~~

// CreateCommand is a command to store a semaphore's context.
type CreateCommand struct {
	*BaseCommand
	CmdName  string
	Capacity int
}

// Do creates file to store a semaphore's context.
func (c *CreateCommand) Do() error {
	var err error

	args := c.FlagSet().Args()
	capacity := c.Capacity
	if len(args) > 0 {
		if capacity, err = strconv.Atoi(args[0]); err != nil {
			return err
		}
	}

	file, err := os.Create(c.BaseCommand.FileName)
	if err != nil {
		return err
	}

	task := Task{Capacity: capacity}
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

// FlagSet returns configured FlagSet to handle CreateCommand arguments.
func (c *CreateCommand) FlagSet() *flag.FlagSet {
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
	}
	return c.Flags
}

// Name ...
func (c *CreateCommand) Name() string {
	return c.CmdName
}

// Desc ...
func (c *CreateCommand) Desc() string {
	return "desc..."
}

// AddCommand is a command to add a job into a semaphore's context.
type AddCommand struct {
	*BaseCommand
	CmdName string
	Command []string
}

// Do adds a job into a semaphore's context.
func (c *AddCommand) Do() error {
	args := c.FlagSet().Args()
	if len(args) == 0 {
		return fmt.Errorf("need args: help call...")
	}

	file, err := os.OpenFile(c.BaseCommand.FileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	var task Task
	data, err := ioutil.ReadAll(file)
	if err := json.Unmarshal(data, &task); err != nil {
		return err
	}

	var jobArgs []string
	if len(args) > 1 {
		jobArgs = args[1:]
	}
	task.AddJob(Job{Name: args[0], Args: jobArgs})

	data, err = json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = file.WriteAt(data, 0)
	return err
}

// FlagSet returns configured FlagSet to handle AddCommand arguments.
func (c *AddCommand) FlagSet() *flag.FlagSet {
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
	}
	return c.Flags
}

// Name ...
func (c *AddCommand) Name() string {
	return c.CmdName
}

// Desc ...
func (c *AddCommand) Desc() string {
	return "desc..."
}

// WaitCommand is a command to execute a semaphore's task.
type WaitCommand struct {
	*BaseCommand
	CmdName        string
	Stdout, Stderr io.Writer
	Notify         bool
	Timeout        time.Duration
}

// Do executes a semaphore's task.
func (c *WaitCommand) Do() error {
	var err error

	file, err := os.OpenFile(c.BaseCommand.FileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	var task Task
	data, err := ioutil.ReadAll(file)
	if err := json.Unmarshal(data, &task); err != nil {
		return err
	}
	if c.Timeout > 0 {
		task.Timeout = c.Timeout
	}

	results := task.Run()
	for result := range results {
		var (
			src io.Reader = result.Stdout
			dst io.Writer = c.Stdout
		)

		if result.Error != nil {
			src, dst = result.Stderr, c.Stderr
		}

		fmt.Fprintf(dst, "command %s: `%s %s`\n", result.Job.ID, result.Job.Name, strings.Join(result.Job.Args, " "))
		fmt.Fprintf(dst, "     error: %v\n", result.Error)
		fmt.Fprint(dst, "    output:\n<<<\n")
		io.Copy(dst, src)
		dst.Write([]byte("\n>>>\n"))
	}

	return err
}

// FlagSet returns configured FlagSet to handle WaitCommand arguments.
func (c *WaitCommand) FlagSet() *flag.FlagSet {
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
		c.Flags.BoolVar(&c.Notify, "notify", false, "")
		c.Flags.DurationVar(&c.Timeout, "timeout", time.Minute, "")
	}
	return c.Flags
}

// Name ...
func (c *WaitCommand) Name() string {
	return c.CmdName
}

// Desc ...
func (c *WaitCommand) Desc() string {
	return "desc..."
}

// HelpCommand ...
type HelpCommand struct {
	*BaseCommand
	CmdName               string
	Commit, Date, Version string
	Commands              Commands
	Error                 error
	Output                io.Writer
}

// FlagSet ...
func (c *HelpCommand) FlagSet() *flag.FlagSet {
	if c.Flags == nil {
		c.Flags = c.BaseCommand.FlagSet(c.CmdName)
	}
	return c.Flags
}

// Name ...
func (c *HelpCommand) Name() string {
	return c.CmdName
}

// Desc ...
func (c *HelpCommand) Desc() string {
	return ""
}

// Do ...
func (c *HelpCommand) Do() error {
	switch c.Error {
	case NotProvided:
		fallthrough
	case flag.ErrHelp:
		c.Usage()
		return nil
	case NotFound:
		c.Usage()
		return c.Error
	default:
		fmt.Fprint(c.Output, c.Error)
		return c.Error
	}
}

func (c *HelpCommand) Usage() {
	fmt.Fprintf(c.Output, `
Usage: %s COMMAND

Semaphore provides functionality to execute terminal commands in parallel.

`, c.BinName)

	if len(c.Commands) > 0 {
		fmt.Fprintln(c.Output, "Commands:")
		for _, cmd := range c.Commands {
			fmt.Fprintf(c.Output, "\n%s\t%s\n", cmd.Name(), cmd.Desc())
			fs := cmd.FlagSet()
			fs.SetOutput(c.Output)
			fs.PrintDefaults()
			fmt.Fprintln(c.Output)
		}
	}

	fmt.Fprintf(c.Output, "Version %s (commit: %s, build date: %s)\n", c.Version, c.Commit, c.Date)
}
