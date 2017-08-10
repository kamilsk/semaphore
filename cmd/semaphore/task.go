package main

import (
	"bytes"
	"io"
	"math"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kamilsk/semaphore"
)

// Task holds required jobs for execution.
type Task struct {
	Capacity int
	Timeout  time.Duration
	Jobs     []Job

	Results []Result `json:"-"`
}

// AddJob adds a job to the task.
func (t *Task) AddJob(job Job) {
	t.Jobs = append(t.Jobs, job)
}

// Run executes all jobs.
func (t *Task) Run() {
	wg := &sync.WaitGroup{}
	sem := semaphore.New(t.Capacity)
	timeout := semaphore.WithTimeout(t.Timeout)

	t.Results = make([]Result, len(t.Jobs))

	var index int32 = math.MaxInt32
	for i := range t.Jobs {
		wg.Add(1)
		go func() {
			result := Result{
				Job:    t.Jobs[i],
				Stdout: bytes.NewBuffer(make([]byte, 1024)),
				Stderr: bytes.NewBuffer(make([]byte, 1024)),
			}

			defer func() {
				t.Results[atomic.AddInt32(&index, 1)] = result
				wg.Done()
			}()

			releaser, err := sem.Acquire(timeout)
			if err != nil {
				result.Error = err
				return
			}
			defer releaser()

			if err := result.Fetch(); err != nil {
				result.Error = err
				return
			}
		}()
	}
	wg.Wait()
}

// Job represents command for execution.
type Job struct {
	ID   string
	Name string
	Args []string
}

// Run executes command.
func (j Job) Run(stdout, stderr io.Writer) error {
	c := exec.Command(j.Name, j.Args...)
	c.Stdout, c.Stderr = stdout, stderr
	return c.Run()
}

// Result holds the job execution result.
type Result struct {
	Error          error
	Job            Job
	Stdout, Stderr *bytes.Buffer
}

// Fetch executes the job and fetches its result into buffers.
func (r Result) Fetch() error {
	return r.Job.Run(r.Stdout, r.Stderr)
}
