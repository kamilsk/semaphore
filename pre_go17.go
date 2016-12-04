// +build !go1.7

package semaphore

import "golang.org/x/net/context"

// Semaphore defines the base interface.
type Semaphore interface {
	HealthChecker
	Releaser

	// Acquire tries to take an available place with the given timeout.
	// If the timeout has occurred, then returns an appropriate error.
	// It must be safe to call Acquire concurrently on a single semaphore.
	Acquire(ctx context.Context) (ReleaseFunc, error)
}

func (sem semaphore) Acquire(ctx context.Context) (ReleaseFunc, error) {
	select {
	case sem <- struct{}{}:
		return releaser(sem), nil
	case <-ctx.Done():
		return nothing, errTimeout
	}
}
