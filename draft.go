package semaphore

import "context"

// Interface is a draft.
type Interface interface {
	HealthChecker

	Acquire(BreakCloser, ...int) (Releaser, error)
	AcquireContext(context.Context, ...int) (Releaser, error)
	Signal(Breaker) <-chan Releaser
}

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
