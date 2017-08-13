package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Command ...
type Command interface {
	Do() error
	FlagSet() *flag.FlagSet
	Name() string
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
		if command == c.Name() {
			var err error
			if len(args) > 1 {
				err = c.FlagSet().Parse(args[1:])
			}
			return c, err
		}
	}
	return nil, errors.New("command not found: help call...")
}

// BaseCommand ...
type BaseCommand struct {
	ID       string
	Mode     flag.ErrorHandling
	Filename string
	fs       *flag.FlagSet
}

// FlagSet ...
func (c *BaseCommand) FlagSet() *flag.FlagSet {
	if c.fs == nil {
		c.fs = flag.NewFlagSet(c.ID, c.Mode)
	}
	return c.fs
}

// Name ...
func (c *BaseCommand) Name() string {
	return c.ID
}

// CreateCommand is a command to store a semaphore's context.
type CreateCommand struct {
	BaseCommand
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

	file, err := os.Create(c.Filename)
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
		c.fs = c.BaseCommand.FlagSet()
	}
	return c.fs
}

// AddCommand is a command to add a job into a semaphore's context.
type AddCommand struct {
	BaseCommand
	Command []string
}

// Do adds a job into a semaphore's context.
func (c *AddCommand) Do() error {
	args := c.FlagSet().Args()
	if len(args) == 0 {
		return errors.New("need args: help call...")
	}

	file, err := os.OpenFile(c.Filename, os.O_RDWR, 0644)
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
		c.fs = c.BaseCommand.FlagSet()
	}
	return c.fs
}

// WaitCommand is a command to execute a semaphore's task.
type WaitCommand struct {
	BaseCommand
	Stdout, Stderr io.Writer
	Notify         bool
	Timeout        time.Duration
}

// Do executes a semaphore's task.
func (c *WaitCommand) Do() error {
	var err error

	file, err := os.OpenFile(c.Filename, os.O_RDWR, 0644)
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

		// TODO combine error
		_, err = fmt.Fprintf(dst, "command %s: `%s %s`\n", result.Job.ID, result.Job.Name, strings.Join(result.Job.Args, " "))
		_, err = fmt.Fprintf(dst, "     error: %v\n", result.Error)
		_, err = fmt.Fprint(dst, "    output:\n<<<\n")
		_, err = io.Copy(dst, src)
		_, err = dst.Write([]byte("\n>>>\n"))
	}

	return err
}

// FlagSet returns configured FlagSet to handle WaitCommand arguments.
func (c *WaitCommand) FlagSet() *flag.FlagSet {
	if c.fs == nil {
		c.fs = c.BaseCommand.FlagSet()
		c.fs.BoolVar(&c.Notify, "notify", false, "")
		c.fs.DurationVar(&c.Timeout, "timeout", time.Minute, "")
	}
	return c.fs
}
