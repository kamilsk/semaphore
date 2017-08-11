package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
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
		return nil, errors.New("help call...")
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
	return nil, errors.New("help call...")
}

// BaseCommand ...
type BaseCommand struct {
	ID   string
	Mode flag.ErrorHandling
	fs   *flag.FlagSet
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
	Filename string
}

// Do creates file to store a semaphore's context.
func (c *CreateCommand) Do() error {
	var err error

	args := c.FlagSet().Args()
	capacity := runtime.GOMAXPROCS(0)
	if len(args) > 0 {
		if capacity, err = strconv.Atoi(args[0]); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(c.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

// FlagSet returns configured FlagSet to handle passed arguments.
func (c *CreateCommand) FlagSet() *flag.FlagSet {
	if c.fs == nil {
		c.fs = c.BaseCommand.FlagSet()
	}
	return c.fs
}

// AddCommand ...
type AddCommand struct {
	BaseCommand
	Command []string
}

// Do ...
func (c *AddCommand) Do() error {
	fmt.Println(c.ID, "run", c.FlagSet().Args(), c.Command)
	return nil
}

// FlagSet ...
func (c *AddCommand) FlagSet() *flag.FlagSet {
	if c.fs == nil {
		c.fs = c.BaseCommand.FlagSet()
	}
	return c.fs
}

// WaitCommand ...
type WaitCommand struct {
	BaseCommand
	Notify  bool
	Timeout time.Duration
}

// Do ...
func (c *WaitCommand) Do() error {
	fmt.Println(c.ID, "run", c.FlagSet().Args(), c.Notify, c.Timeout)
	return nil
}

// FlagSet ...
func (c *WaitCommand) FlagSet() *flag.FlagSet {
	if c.fs == nil {
		c.fs = c.BaseCommand.FlagSet()
		c.fs.BoolVar(&c.Notify, "notify", false, "")
		c.fs.DurationVar(&c.Timeout, "timeout", time.Minute, "")
	}
	return c.fs
}
