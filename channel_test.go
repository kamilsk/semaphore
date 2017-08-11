package semaphore_test

import (
	"os"
	"testing"
	"time"

	"github.com/kamilsk/semaphore"
)

func TestWithTimeout(t *testing.T) {
	sleep := 500 * time.Millisecond

	start := time.Now()
	<-semaphore.WithTimeout(sleep)
	end := time.Now()

	if expected, obtained := sleep, end.Sub(start); expected > obtained {
		t.Errorf("unexpected sleep time. expected: %s; obtained: %s", expected, obtained)
	}
}

func TestWithSignal_NilSignal(t *testing.T) {
	<-semaphore.WithSignal(nil)
}

func TestMultiplex(t *testing.T) {
	sleep := 500 * time.Millisecond

	start := time.Now()
	<-semaphore.Multiplex(semaphore.WithSignal(os.Interrupt), semaphore.WithTimeout(sleep))
	end := time.Now()

	if expected, obtained := sleep, end.Sub(start); expected > obtained {
		t.Errorf("unexpected sleep time. expected: %s; obtained: %s", expected, obtained)
	}
}

func TestMultiplex_WithoutChannels(t *testing.T) {
	<-semaphore.Multiplex()
}
