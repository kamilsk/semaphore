package semaphore

import (
	"runtime"
	"testing"
	"time"
)

func TestCapacity(t *testing.T) {
	if obtained, expected := Capacity(), runtime.GOMAXPROCS(0); obtained != expected {
		t.Errorf("an unexpected capacity. expected: %d; obtained: %d", expected, obtained)
	}
}

func TestOccupied(t *testing.T) {
	if obtained, expected := Occupied(), 0; obtained != expected {
		t.Errorf("unexpected occupied places. expected: %d; obtained: %d", expected, obtained)
	}
}

func TestRelease(t *testing.T) {
	if err, expected := Release(), "semaphore is empty"; err.Error() != expected {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
}

func TestAcquire(t *testing.T) {
	rs := make([]ReleaseFunc, 0, Capacity())
	do := func() {
		for _, r := range rs {
			r()
		}
	}
	for i := 0; i < Capacity(); i++ {
		r, _ := Acquire(nil)
		rs = append(rs, r)
	}
	expected := "operation timeout"
	if _, err := Acquire(WithTimeout(10 * time.Millisecond)); err.Error() != expected {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
	do()
	if r, err := Acquire(WithTimeout(10 * time.Millisecond)); err != nil {
		t.Error("an unexpected error", err)
	} else {
		r()
	}
}
