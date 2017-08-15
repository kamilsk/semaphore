package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Command ...
type Command interface {
	Do() error
	FlagSet() *flag.FlagSet
	Name() string
	Desc() string
}

// Commands ...
type Commands []Command

// Parse ...
func (l Commands) Parse(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, errors.New("need a command: help call...")
	}
	command := args[0]
	for _, c := range l {
		if c.Name() == command {
			var err error
			if len(args) > 1 {
				err = c.FlagSet().Parse(args[1:])
			}
			return c, err
		}
	}
	return nil, errors.New("command not found: help call...")
}

// ~~~

// BaseCommand ...
type BaseCommand struct {
	BinName  string
	FileName string
	Mode     flag.ErrorHandling

	fs *flag.FlagSet
}

// Copy ...
func (c *BaseCommand) Copy() *BaseCommand {
	n := *c
	return &n
}

// FlagSet ...
func (c *BaseCommand) FlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, c.Mode)
	fs.StringVar(&c.FileName, "filename", filepath.Join(os.TempDir(), c.BinName+".json"), "")
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
	if c.fs == nil {
		c.fs = c.BaseCommand.FlagSet(c.CmdName)
	}
	return c.fs
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
		return errors.New("need args: help call...")
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
	if c.fs == nil {
		c.fs = c.BaseCommand.FlagSet(c.CmdName)
	}
	return c.fs
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
	if c.fs == nil {
		c.fs = c.BaseCommand.FlagSet(c.CmdName)
		c.fs.BoolVar(&c.Notify, "notify", false, "")
		c.fs.DurationVar(&c.Timeout, "timeout", time.Minute, "")
	}
	return c.fs
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
	Commands              []Command
	Output                io.Writer
}

// Do ...
func (c *HelpCommand) Do() error {
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

	return nil
}

// FlagSet ...
func (c *HelpCommand) FlagSet() *flag.FlagSet {
	if c.fs == nil {
		c.fs = c.BaseCommand.FlagSet(c.CmdName)
	}
	return c.fs
}

// Name ...
func (c *HelpCommand) Name() string {
	return c.CmdName
}
