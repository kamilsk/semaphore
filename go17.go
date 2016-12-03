// +build go1.7

package semaphore

import "context"

// Semaphore defines the base interface.
type Semaphore interface {
	HealthChecker
	Releaser

	// Acquire tries to take an available place with the given context.
	// If the context will done, then returns an appropriate error.
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
