package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/kamilsk/semaphore"
)

// Task holds required jobs for execution.
type Task struct {
	Capacity int
	Timeout  time.Duration
	Jobs     []Job
}

// AddJob adds a job to the task.
func (t *Task) AddJob(job Job) {
	job.ID = fmt.Sprintf("#%d", len(t.Jobs)+1)
	t.Jobs = append(t.Jobs, job)
}

// Run executes all jobs.
func (t *Task) Run() <-chan Result {
	results := make(chan Result, len(t.Jobs))

	go func() {
		defer func() { close(results) }()

		wg := &sync.WaitGroup{}
		sem := semaphore.New(t.Capacity)
		deadline := semaphore.Multiplex(
			semaphore.WithTimeout(t.Timeout),
			semaphore.WithSignal(os.Interrupt),
		)

		for i := range t.Jobs {
			wg.Add(1)
			go func(index int) {
				result := Result{
					Job:    t.Jobs[index],
					Stdout: bytes.NewBuffer(make([]byte, 1024)),
					Stderr: bytes.NewBuffer(make([]byte, 1024)),
				}

				defer func() {
					results <- result
					wg.Done()
				}()

				releaser, err := sem.Acquire(deadline)
				if err != nil {
					result.Error = err
					return
				}
				defer releaser()

				if err := result.Fetch(); err != nil {
					result.Error = err
					return
				}
			}(i)
		}
		wg.Wait()
	}()

	return results
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
	Job            Job
	Error          error
	Stdout, Stderr *bytes.Buffer
}

// Fetch executes the job and fetches its result into buffers.
func (r Result) Fetch() error {
	return r.Job.Run(r.Stdout, r.Stderr)
}
