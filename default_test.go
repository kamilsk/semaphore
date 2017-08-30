package semaphore_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/kamilsk/semaphore"
)

func TestAcquire(t *testing.T) {
	rs := make([]semaphore.ReleaseFunc, 0, semaphore.Capacity())
	do := func() {
		for _, r := range rs {
			r()
		}
	}
	for i := 0; i < semaphore.Capacity(); i++ {
		r, _ := semaphore.Acquire(nil)
		rs = append(rs, r)
	}
	expected := "operation timeout"
	if _, err := semaphore.Acquire(semaphore.WithTimeout(10 * time.Millisecond)); err.Error() != expected {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
	do()
	if r, err := semaphore.Acquire(semaphore.WithTimeout(10 * time.Millisecond)); err != nil {
		t.Error("an unexpected error", err)
	} else {
		r()
	}
}

func TestCapacity(t *testing.T) {
	if obtained, expected := semaphore.Capacity(), runtime.GOMAXPROCS(0); obtained != expected {
		t.Errorf("an unexpected capacity. expected: %d; obtained: %d", expected, obtained)
	}
}

func TestOccupied(t *testing.T) {
	if obtained, expected := semaphore.Occupied(), 0; obtained != expected {
		t.Errorf("unexpected occupied places. expected: %d; obtained: %d", expected, obtained)
	}
}

func TestRelease(t *testing.T) {
	if err, expected := semaphore.Release(), "semaphore is empty"; err.Error() != expected {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
}

func TestSignal(t *testing.T) {
	if release, ok := <-semaphore.Signal(nil); release == nil || !ok {
		t.Error("unexpected signal")
	}
}
