package semaphore

import (
	"errors"
	"time"
)

// Semaphore defines the base interface.
type Semaphore interface {
	HealthCheck
	// Acquire tries to take an available place with the given timeout.
	// If the timeout has occurred, then returns an appropriate error.
	// It must be safe to call Acquire concurrently on a single semaphore.
	Acquire(time.Duration) error
	// Release releases the previously occupied place.
	// If no places was occupied then returns error.
	// It must be safe to call Release concurrently on a single semaphore.
	Release() error
}

// HealthCheck defines some helpful methods related with capacity for monitoring.
type HealthCheck interface {
	// Capacity returns the number of total places.
	// It must be safe to call Capacity concurrently on a single semaphore.
	Capacity() int
	// Occupied returns the number of occupied places.
	// It must be safe to call Occupied concurrently on a single semaphore.
	Occupied() int
}

// New constructs a new thread-safe Semaphore with the given number of places.
func New(size int) Semaphore {
	return make(semaphore, size)
}

var (
	errEmpty   = errors.New("semaphore is empty")
	errTimeout = errors.New("operation timeout")
)

type semaphore chan struct{}

func (sem semaphore) Acquire(timeout time.Duration) error {
	// returns errTimeout immediately
	// without unnecessary overhead
	if timeout < 0 {
		return errTimeout
	}
	select {
	case sem <- struct{}{}:
		return nil
	case <-time.After(timeout):
		return errTimeout
	}
}

func (sem semaphore) Release() error {
	select {
	case _ = <-sem:
		return nil
	default:
		return errEmpty
	}
}

func (sem semaphore) Capacity() int {
	return cap(sem)
}

func (sem semaphore) Occupied() int {
	return len(sem)
}
