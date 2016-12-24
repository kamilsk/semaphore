package simple

import (
	"errors"
	"time"

	intf "github.com/kamilsk/semaphore"
)

// Semaphore provides the functionality of the same named pattern.
type Semaphore interface {
	intf.HealthChecker
	intf.Releaser

	// Acquire tries to take an available occupancy with the given timeout.
	// If the timeout has occurred, then returns an appropriate error.
	// It must be safe to call Acquire concurrently on a single semaphore.
	Acquire(timeout time.Duration) error
}

// New constructs a new thread-safe Semaphore with the given capacity.
func New(capacity int) Semaphore {
	return make(semaphore, capacity)
}

var (
	errEmpty   = errors.New("semaphore is empty")
	errTimeout = errors.New("operation timeout")
)

type semaphore chan struct{}

func (sem semaphore) Acquire(timeout time.Duration) error {
	if timeout <= 0 {
		return errTimeout
	}
	select {
	case sem <- struct{}{}:
		return nil
	case <-time.After(timeout):
		return errTimeout
	}
}

func (sem semaphore) Capacity() int {
	return cap(sem)
}

func (sem semaphore) Occupied() int {
	return len(sem)
}

func (sem semaphore) Release() error {
	select {
	case <-sem:
		return nil
	default:
		return errEmpty
	}
}
