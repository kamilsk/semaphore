package main

import (
	"errors"
	"flag"
	"fmt"
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
			if len(args) > 1 {
				c.FlagSet().Parse(args[1:])
			}
			return c, nil
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

// CreateCommand ...
type CreateCommand struct {
	BaseCommand
	Capacity int
}

// Do ...
func (c *CreateCommand) Do() error {
	fmt.Println(c.ID, "run", c.FlagSet().Args(), c.Capacity)
	return nil
}

// FlagSet ...
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
