package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTask_Success(t *testing.T) {
	task := &Task{Capacity: 2, Timeout: time.Second}
	task.AddJob(Job{ID: "test", Name: "echo", Args: []string{"hello,", "world"}})
	task.AddJob(Job{ID: "test", Name: "echo", Args: []string{"hello,", "world"}})
	r := task.Run()

	assert.Equal(t, cap(r), 2)
	assert.NoError(t, (<-r).Error)
	assert.NoError(t, (<-r).Error)
}

func TestTask_Fails(t *testing.T) {
	{
		task := &Task{}
		task.AddJob(Job{ID: "test", Name: "echo", Args: []string{"hello,", "world"}})
		r := task.Run()

		assert.Equal(t, cap(r), 1)
		assert.Error(t, (<-r).Error)
	}
	{
		task := &Task{Capacity: 1, Timeout: time.Second}
		task.AddJob(Job{ID: "test", Name: "curl", Args: []string{"unknown"}})
		r := task.Run()

		assert.Equal(t, cap(r), 1)
		assert.Error(t, (<-r).Error)
	}
}

func TestJob(t *testing.T) {
	j := Job{ID: "test", Name: "echo", Args: []string{"hello,", "world"}}

	assert.Equal(t, "echo#test", fmt.Sprintf("%v", j))
	assert.Equal(t, "echo#test [hello, world]", fmt.Sprintf("%+v", j))
	assert.Equal(t, "echo#test", j.String())
	assert.Equal(t, "echo#test `echo hello, world`", fmt.Sprintf("%q", j))
	assert.NoError(t, j.Run(ioutil.Discard, ioutil.Discard))
	assert.Equal(t, "echo#test", j.String())
}

func TestResult(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	r := Result{Job: Job{ID: "test", Name: "echo", Args: []string{"hello,", "world"}}, Stderr: buf, Stdout: buf}

	assert.NoError(t, r.Fetch())
}

func TestResults(t *testing.T) {
	r := Results{}
	r.Append(Result{Job: Job{ID: "3"}})
	r.Append(Result{Job: Job{ID: "1"}})
	r.Append(Result{Job: Job{ID: "2"}})

	assert.Equal(t, 3, r.Len())
	assert.Equal(t, Results{Result{Job: Job{ID: "1"}}, Result{Job: Job{ID: "2"}}, Result{Job: Job{ID: "3"}}}, r.Sort())
}
