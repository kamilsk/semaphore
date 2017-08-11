package semaphore

import (
	"os"
	"os/signal"
	"reflect"
	"time"
)

// WithTimeout returns empty struct channel based on Time channel.
func WithTimeout(timeout time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	if timeout <= 0 {
		close(ch)
		return ch
	}
	go func() {
		<-time.After(timeout)
		close(ch)
	}()
	return ch
}

// WithSignal returns empty struct channel based on Signal channel.
func WithSignal(s os.Signal) <-chan struct{} {
	ch := make(chan struct{})
	if s == nil {
		close(ch)
		return ch
	}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, s)
		<-c
		close(ch)
		signal.Stop(c)
	}()
	return ch
}

// Multiplex combines multiple empty struct channels into one.
func Multiplex(channels ...<-chan struct{}) <-chan struct{} {
	ch := make(chan struct{})
	if len(channels) == 0 {
		close(ch)
		return ch
	}
	go func() {
		cases := make([]reflect.SelectCase, 0, len(channels))
		for _, ch := range channels {
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
		}
		reflect.Select(cases)
		close(ch)
	}()
	return ch
}
