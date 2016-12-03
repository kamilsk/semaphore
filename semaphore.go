package semaphore

import "errors"

// HealthChecker defines some helpful methods related with semaphore's state.
type HealthChecker interface {
	// Capacity returns the number of total places.
	// It must be safe to call Capacity concurrently on a single semaphore.
	Capacity() int
	// Occupied returns the number of occupied places.
	// It must be safe to call Occupied concurrently on a single semaphore.
	Occupied() int
}

// Releaser defines method to release the previously occupied place.
type Releaser interface {
	// Release releases the previously occupied place.
	// If no places was occupied then returns error.
	// It must be safe to call Release concurrently on a single semaphore.
	Release() error
}

// A ReleaseFunc tells a semaphore to release the previously occupied place and to ignore an error.
type ReleaseFunc func()

// New constructs a new thread-safe Semaphore with the given number of places.
func New(capacity int) Semaphore {
	return make(semaphore, capacity)
}

var (
	nothing ReleaseFunc = func() {}

	errEmpty   = errors.New("semaphore is empty")
	errTimeout = errors.New("operation timeout")
)

type semaphore chan struct{}

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

func releaser(releaser Releaser) ReleaseFunc {
	return func() { _ = releaser.Release() }
}
