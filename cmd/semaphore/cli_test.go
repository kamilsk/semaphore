package main_test

import (
	"flag"
	"testing"

	"github.com/kamilsk/semaphore/cmd/semaphore"
)

func TestCreateCommand(t *testing.T) {
	bc := &main.BaseCommand{BinName: "test"}
	cc := &main.CreateCommand{BaseCommand: bc, CmdName: "test"}

	if expected, obtained := cc.CmdName, cc.Name(); expected != obtained {
		t.Errorf("unexpected command name. expected: %s, obtained: %s", expected, obtained)
	}

	if cc.Desc() == "" {
		t.Error("unexpected empty description")
	}

	if expected, obtained := 1, func() int {
		var count int
		cc.FlagSet().VisitAll(func(*flag.Flag) { count++ })
		return count
	}(); expected != obtained {
		t.Errorf("unexpected command flags count. expected: %d, obtained: %d", expected, obtained)
	}
}

func TestAddCommand(t *testing.T) {
	bc := &main.BaseCommand{BinName: "test"}
	cc := &main.AddCommand{BaseCommand: bc, CmdName: "test"}

	if expected, obtained := cc.CmdName, cc.Name(); expected != obtained {
		t.Errorf("unexpected command name. expected: %s, obtained: %s", expected, obtained)
	}

	if cc.Desc() == "" {
		t.Error("unexpected empty description")
	}

	if expected, obtained := 2, func() int {
		var count int
		cc.FlagSet().VisitAll(func(*flag.Flag) { count++ })
		return count
	}(); expected != obtained {
		t.Errorf("unexpected command flags count. expected: %d, obtained: %d", expected, obtained)
	}
}

func TestWaitCommand(t *testing.T) {
	bc := &main.BaseCommand{BinName: "test"}
	cc := &main.WaitCommand{BaseCommand: bc, CmdName: "test"}

	if expected, obtained := cc.CmdName, cc.Name(); expected != obtained {
		t.Errorf("unexpected command name. expected: %s, obtained: %s", expected, obtained)
	}

	if cc.Desc() == "" {
		t.Error("unexpected empty description")
	}

	if expected, obtained := 4, func() int {
		var count int
		cc.FlagSet().VisitAll(func(*flag.Flag) { count++ })
		return count
	}(); expected != obtained {
		t.Errorf("unexpected command flags count. expected: %d, obtained: %d", expected, obtained)
	}
}

func TestHelpCommand(t *testing.T) {
	bc := &main.BaseCommand{BinName: "test"}
	cc := &main.HelpCommand{BaseCommand: bc, CmdName: "test"}

	if expected, obtained := cc.CmdName, cc.Name(); expected != obtained {
		t.Errorf("unexpected command name. expected: %s, obtained: %s", expected, obtained)
	}

	if cc.Desc() == "" {
		t.Error("unexpected empty description")
	}

	if expected, obtained := 1, func() int {
		var count int
		cc.FlagSet().VisitAll(func(*flag.Flag) { count++ })
		return count
	}(); expected != obtained {
		t.Errorf("unexpected command flags count. expected: %d, obtained: %d", expected, obtained)
	}
}
