// +build !go1.7

package semaphore

import "golang.org/x/net/context"

// Semaphore provides the functionality of the same named pattern.
type Semaphore interface {
	HealthChecker
	Releaser

	// Acquire tries to reduces the number of available slots for 1.
	// The operation can be canceled using context. In this case
	// an appropriate error will be returned.
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
