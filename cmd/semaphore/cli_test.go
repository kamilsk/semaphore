package main

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	bc := &BaseCommand{BinName: "test"}
	cc := &CreateCommand{BaseCommand: bc, CmdName: "test"}

	assert.Equal(t, cc.CmdName, cc.Name())
	assert.NotEmpty(t, cc.Desc())
	assert.Equal(t, 2, func() int {
		var count int
		cc.FlagSet().VisitAll(func(*flag.Flag) { count++ })
		return count
	}())
}

func TestAddCommand(t *testing.T) {
	bc := &BaseCommand{BinName: "test"}
	cc := &AddCommand{BaseCommand: bc, CmdName: "test"}

	assert.Equal(t, cc.CmdName, cc.Name())
	assert.NotEmpty(t, cc.Desc())
	assert.Equal(t, 3, func() int {
		var count int
		cc.FlagSet().VisitAll(func(*flag.Flag) { count++ })
		return count
	}())
}

func TestWaitCommand(t *testing.T) {
	bc := &BaseCommand{BinName: "test"}
	cc := &WaitCommand{BaseCommand: bc, CmdName: "test"}

	assert.Equal(t, cc.CmdName, cc.Name())
	assert.NotEmpty(t, cc.Desc())
	assert.Equal(t, 5, func() int {
		var count int
		cc.FlagSet().VisitAll(func(*flag.Flag) { count++ })
		return count
	}())
}

func TestHelpCommand(t *testing.T) {
	bc := &BaseCommand{BinName: "test"}
	cc := &HelpCommand{BaseCommand: bc, CmdName: "test"}

	assert.Equal(t, cc.CmdName, cc.Name())
	assert.NotEmpty(t, cc.Desc())
	assert.Equal(t, 2, func() int {
		var count int
		cc.FlagSet().VisitAll(func(*flag.Flag) { count++ })
		return count
	}())
}
