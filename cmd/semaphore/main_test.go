package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain_Exec_Fails(t *testing.T) {
	var status int
	Main{
		Args:   []string{"cmd", "unknown"},
		Stdout: ioutil.Discard, Stderr: ioutil.Discard,
		Shutdown: func(code int) { status = code },
	}.Exec()

	assert.Equal(t, 1, status)
}

func TestMain_Exec__Create(t *testing.T) {
	var status int
	Main{
		Args:   []string{"cmd", "create", "not int"},
		Stdout: ioutil.Discard, Stderr: ioutil.Discard,
		Shutdown: func(code int) { status = code },
	}.Exec()

	assert.Equal(t, 1, status)
}

func TestMain_Exec__Help(t *testing.T) {
	var status int
	{
		Main{
			Args:   []string{"cmd", "help"},
			Stdout: ioutil.Discard, Stderr: ioutil.Discard,
			Shutdown: func(code int) { status = code },
		}.Exec()

		assert.Equal(t, 0, status)
	}
	{
		Main{
			Args:   []string{"cmd"},
			Stdout: ioutil.Discard, Stderr: ioutil.Discard,
			Shutdown: func(code int) { status = code },
		}.Exec()

		assert.Equal(t, 0, status)
	}
}
