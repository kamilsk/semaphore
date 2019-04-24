package semaphore

// A Breaker carries a cancellation signal to break an action execution.
//
// It is a subset of context.Context and github.com/kamilsk/breaker.Breaker.
type Breaker interface {
	// Done returns a channel that's closed when a cancellation signal occurred.
	Done() <-chan struct{}
}

// A BreakCloser carries a cancellation signal to break an action execution
// and can release resources associated with it.
//
// It is a subset of github.com/kamilsk/breaker.Breaker.
type BreakCloser interface {
	Breaker
	// Close closes the Done channel and releases resources associated with it.
	Close()
}

// A Releaser provides a possibility to release resources that it holds.
type Releaser interface {
	// Release releases resources associated with the Releaser.
	// If no places were occupied, then it returns an appropriate error.
	Release() error
}

// Interface defines the functionality of the Semaphore pattern.
type Interface interface {
	Releaser

	Acquire(BreakCloser, ...uint32) (Releaser, error)
	Try(Breaker, ...uint32) (Releaser, error)
	Signal(Breaker) <-chan Releaser

	Peek() uint32
	Size(uint) uint32
}

// Semaphore provides the functionality of the same named pattern.
// Deprecated: will be replaced by Interface.
type Semaphore interface {
	HealthChecker
	Releaser

	// Acquire tries to reduce the number of available slots for 1.
	// The operation can be canceled using context. In this case,
	// it returns an appropriate error.
	Acquire(deadline <-chan struct{}) (ReleaseFunc, error)
	// Catch tries to reduce the number of available slots for 1.
	Catch() (ReleaseFunc, error)
	// Signal returns a channel to send to it release function
	// only if Acquire is successful. In any case, the channel will be closed.
	Signal(deadline <-chan struct{}) <-chan ReleaseFunc
}

// HealthChecker defines helpful methods related with semaphore status.
// Deprecated: will be replaced by Interface.
type HealthChecker interface {
	// Capacity returns a capacity of a semaphore.
	// It must be safe to call Capacity concurrently on a single semaphore.
	Capacity() int
	// Occupied returns a current number of occupied slots.
	// It must be safe to call Occupied concurrently on a single semaphore.
	Occupied() int
}
