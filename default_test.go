package semaphore

import (
	"context"
	"runtime"
	"testing"
	"time"
)

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	expected := "operation timeout"
	if _, err := Acquire(ctx.Done()); err.Error() != expected {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
	do()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
	if r, err := Acquire(ctx.Done()); err != nil {
		t.Error("an unexpected error", err)
	} else {
		r()
	}
}

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

func TestSignal(t *testing.T) {
	if release, ok := <-Signal(nil); release == nil || !ok {
		t.Error("unexpected signal")
	}
}
