// +build go1.7

package semaphore

import "context"

// WithContext returns Context with cancellation based on empty struct channel.
func WithContext(deadline <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-deadline
		cancel()
		return
	}()
	return ctx
}
