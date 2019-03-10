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

// IsTimeout checks if passed error is related to call Acquire on full semaphore.
func IsTimeout(err error) bool {
	return err == errTimeout
}

var (
	nothing ReleaseFunc = func() {}

	errEmpty   = errors.New("semaphore is empty")
	errTimeout = errors.New("operation timeout")
)

type semaphore chan struct{}

func (sem semaphore) Acquire(deadline <-chan struct{}) (ReleaseFunc, error) {
	select {
	case sem <- struct{}{}:
		return func() { _ = sem.Release() }, nil //nolint: gas
	case <-deadline:
		return nothing, errTimeout
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

func (sem semaphore) Signal(deadline <-chan struct{}) <-chan ReleaseFunc {
	ch := make(chan ReleaseFunc, 1)
	go func() {
		if release, err := sem.Acquire(deadline); err == nil {
			ch <- release
		}
		close(ch)
	}()
	return ch
}
