package semaphore // import "github.com/kamilsk/semaphore"

import "errors"

// HealthChecker defines helpful methods related with semaphore status.
type HealthChecker interface {
	// Capacity returns the capacity of semaphore.
	// It must be safe to call Capacity concurrently on a single semaphore.
	Capacity() int
	// Occupied returns the current number of occupied slots.
	// It must be safe to call Occupied concurrently on a single semaphore.
	Occupied() int
}

// Releaser defines method to release the previously occupied semaphore.
type Releaser interface {
	// Release releases the previously occupied slot.
	// If no places was occupied then returns an appropriate error.
	// It must be safe to call Release concurrently on a single semaphore.
	Release() error
}

// A ReleaseFunc tells a semaphore to release the previously occupied slot
// and to ignore an error if it occur.
type ReleaseFunc func()

// New constructs a new thread-safe Semaphore with the given capacity.
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
