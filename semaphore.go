// Package semaphore provides an implementation of Semaphore pattern
// with timeout of lock/unlock operations based on channels.
package semaphore

import "errors"

// ReleaseFunc tells a semaphore to release the previously occupied slot
// and ignore an error if it occurs.
type ReleaseFunc func()

// Release calls f().
func (f ReleaseFunc) Release() error {
	f()
	return nil
}

// New constructs a new thread-safe Semaphore with the given capacity.
func New(capacity int) Semaphore {
	return make(semaphore, capacity)
}

// IsEmpty checks if passed error is related to call Release on empty semaphore.
func IsEmpty(err error) bool {
	return err == errEmpty
}

// IsNoPlace checks if passed error is related to call Catch on full semaphore.
func IsNoPlace(err error) bool {
	return err == errNoPlace
}

// IsTimeout checks if passed error is related to call Acquire on full semaphore.
func IsTimeout(err error) bool {
	return err == errTimeout
}

var (
	nothing ReleaseFunc = func() {}

	errEmpty   = errors.New("semaphore is empty")
	errNoPlace = errors.New("semaphore has no place")
	errTimeout = errors.New("operation timeout")
)

type semaphore chan struct{}

func (semaphore semaphore) Acquire(deadline <-chan struct{}) (ReleaseFunc, error) {
	select {
	case semaphore <- struct{}{}:
		return func() { _ = semaphore.Release() }, nil //nolint: gas
	case <-deadline:
		return nothing, errTimeout
	}
}

func (semaphore semaphore) Catch() (ReleaseFunc, error) {
	select {
	case semaphore <- struct{}{}:
		return func() { _ = semaphore.Release() }, nil //nolint: gas
	default:
		return nothing, errNoPlace
	}
}

func (semaphore semaphore) Capacity() int {
	return cap(semaphore)
}

func (semaphore semaphore) Occupied() int {
	return len(semaphore)
}

func (semaphore semaphore) Release() error {
	select {
	case <-semaphore:
		return nil
	default:
		return errEmpty
	}
}

func (semaphore semaphore) Signal(deadline <-chan struct{}) <-chan ReleaseFunc {
	ch := make(chan ReleaseFunc, 1)
	go func() {
		if release, err := semaphore.Acquire(deadline); err == nil {
			ch <- release
		}
		close(ch)
	}()
	return ch
}
