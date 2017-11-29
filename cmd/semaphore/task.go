package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kamilsk/semaphore"
	"github.com/pkg/errors"
)

// Task holds required jobs for execution.
type Task struct {
	Capacity int
	Timeout  time.Duration

	jobs []Job
}

// AddJob sets ID to a job and adds it to the task.
func (t *Task) AddJob(job Job) {
	job.ID = strconv.Itoa(len(t.jobs) + 1)
	t.jobs = append(t.jobs, job)
}

// Run executes all jobs.
func (t *Task) Run() <-chan Result {
	results := make(chan Result, len(t.jobs))

	go func() {
		defer func() { close(results) }()

		sem := semaphore.New(t.Capacity)
		deadline := semaphore.Multiplex(
			semaphore.WithTimeout(t.Timeout),
			semaphore.WithSignal(os.Interrupt),
		)

		wg := &sync.WaitGroup{}
		for i := range t.jobs {
			wg.Add(1)
			go func(job Job) {
				result := Result{
					Job:    job,
					Stdout: bytes.NewBuffer(make([]byte, 1024)),
					Stderr: bytes.NewBuffer(make([]byte, 1024)),
				}

				defer func() {
					result.End = time.Now()
					results <- result
					wg.Done()
				}()

				release, err := sem.Acquire(deadline)
				if err != nil {
					result.Error = errors.WithMessage(err, "semaphore")
					return
				}
				defer release()

				result.Start = time.Now()
				if err := result.Fetch(); err != nil {
					result.Error = err
					return
				}
			}(t.jobs[i])
		}
		wg.Wait()
	}()

	return results
}

// Job represents a command for execution.
type Job struct {
	ID   string
	Name string
	Args []string
}

// Format implements `fmt.Formatter` interface.
func (j Job) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') && len(j.Args) > 0 {
			fmt.Fprintf(s, "%s %+v", j.String(), j.Args)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, j.String())
	case 'q':
		fmt.Fprintf(s, "%s `%s`", j.String(), strings.Join(append([]string{j.Name}, j.Args...), " "))
	}
}

// Run prepares command and executes it.
func (j Job) Run(stdout, stderr io.Writer) error {
	c := exec.Command(j.Name, j.Args...)
	c.Stdout, c.Stderr = stdout, stderr
	return errors.WithMessage(c.Run(), fmt.Sprintf("an error occurred while executing %q", j))
}

// String implements `fmt.Stringer` interface.
func (j Job) String() string {
	return j.Name + "#" + j.ID
}

// Result holds the job execution result.
type Result struct {
	Job            Job
	Error          error
	Stdout, Stderr io.ReadWriter
	Start, End     time.Time
}

// Fetch executes the job and fetches its result into buffers.
func (r Result) Fetch() error {
	return errors.WithMessage(r.Job.Run(r.Stdout, r.Stderr), fmt.Sprintf("the job %s ended with an error", r.Job))
}

// Results is a container implements `sort.Interface`.
type Results []Result

// Append adds result into a container.
func (l *Results) Append(r Result) {
	*l = append(*l, r)
}

// Len returns a container size.
func (l Results) Len() int {
	return len(l)
}

// Less compares two results from container with indexes i and j.
func (l Results) Less(i, j int) bool {
	return strings.Compare(l[i].Job.ID, l[j].Job.ID) == -1
}

// Swap swaps two results from container with indexes i and j.
func (l Results) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Sort sorts results by its ID.
func (l Results) Sort() Results {
	sort.Sort(l)
	return l
}
