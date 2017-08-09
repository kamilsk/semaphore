package semaphore_test

import (
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
