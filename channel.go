package semaphore

import "time"

// WithTimeout returns empty struct channel based on Time channel.
func WithTimeout(timeout time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	if timeout <= 0 {
		close(ch)
		return ch
	}
	go func() {
		for range time.After(timeout) {
			close(ch)
			return
		}
	}()
	return ch
}
