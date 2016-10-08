package semaphore

import (
	"errors"
	"time"
)

// Semaphore defines the base interface.
type Semaphore interface {
	// Acquire tries to take an available place with the given timeout.
	// If the timeout has occurred, then returns an appropriate error.
	Acquire(time.Duration) error
	// Release releases the previously occupied place.
	Release()
}

// New constructs a new Semaphore with the given number of places.
func New(size int) Semaphore {
	return make(semaphore, size)
}

var errTimeout = errors.New("operation timeout")

type semaphore chan struct{}

func (sem semaphore) Acquire(timeout time.Duration) error {
	select {
	case sem <- struct{}{}:
		return nil
	case <-time.After(timeout):
		return errTimeout
	}
}

func (sem semaphore) Release() {
	<-sem
}
