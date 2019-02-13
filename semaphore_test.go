package semaphore_test

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/kamilsk/semaphore/v4"
)

func TestSemaphore_Acquire_Timeout(t *testing.T) {
	expected := "operation timeout"
	for _, tc := range []struct {
		name    string
		timeout time.Duration
	}{
		{name: "negative timeout", timeout: -time.Second},
		{name: "zero timeout", timeout: 0},
		{name: "positive timeout", timeout: time.Nanosecond},
	} {
		sem := New(0)
		release, err := sem.Acquire(WithTimeout(tc.timeout))
		if !IsTimeout(err) {
			t.Errorf("an unexpected error in test case %q. expected: %s; obtained: %v", tc.name, expected, err)
		}
		_ = release.Release()
	}
}

func TestSemaphore_Capacity_Immutability(t *testing.T) {
	capacity := 7

	sem := New(capacity)

	if sem.Capacity() != capacity {
		t.Errorf("an unexpected capacity. expected: %d; obtained: %d", capacity, sem.Capacity())
	}

	for i := 0; i < sem.Capacity(); i++ {
		_, _ = sem.Acquire(nil)
	}

	if sem.Capacity() != capacity {
		t.Errorf("an unexpected capacity. expected: %d; obtained: %d", capacity, sem.Capacity())
	}
}

func TestSemaphore_Occupied_Linearity(t *testing.T) {
	sem := New(7)

	for i := 0; i < sem.Capacity(); i++ {
		if sem.Occupied() != i {
			t.Errorf("unexpected occupied places. expected: %d; obtained: %d", i, sem.Occupied())
		}
		_, _ = sem.Acquire(nil)
	}

	if sem.Occupied() != sem.Capacity() {
		t.Errorf("unexpected occupied places. expected: %d; obtained: %d", sem.Capacity(), sem.Occupied())
	}
}

func TestSemaphore_Release_TryToGetDeadLock(t *testing.T) {
	sem := New(0)

	if err, expected := sem.Release(), "semaphore is empty"; !IsEmpty(err) {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
}

func TestSemaphore_Signal(t *testing.T) {
	sem := New(0)

	release, ok := <-sem.Signal(WithTimeout(0))
	if release != nil || ok {
		t.Error("unexpected signal")
	}
}

func TestSemaphore_Concurrently(t *testing.T) {
	sem := New(int(math.Max(2.0, float64(runtime.GOMAXPROCS(0)))))

	var counter int32

	start, wg := make(chan bool), &sync.WaitGroup{}
	for i := 0; i < sem.Capacity(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			release, err := sem.Acquire(nil)
			if err != nil {
				t.Error("an unexpected error", err)
				return
			}
			defer release()
			atomic.AddInt32(&counter, 1)
		}()
	}
	close(start)
	wg.Wait()

	if int(counter) != sem.Capacity() {
		t.Errorf("an unexpected counter value. expected: %d; obtained: %d", sem.Capacity(), counter)
	}

	if sem.Occupied() != 0 {
		t.Errorf("zero occupied places are expected but received %d instead", sem.Occupied())
	}
}
