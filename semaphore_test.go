package semaphore

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func (sem semaphore) Flush() {
	close(sem)
	for range sem {
	}
}

func TestSemaphore_Acquire_Timeout(t *testing.T) {
	for _, test := range []struct {
		name    string
		timeout time.Duration
	}{
		{name: "negative timeout", timeout: -time.Second},
		{name: "zero timeout", timeout: 0},
		{name: "positive timeout", timeout: time.Nanosecond},
	} {
		sem := New(0)
		release, err := sem.Acquire(WithTimeout(test.timeout))
		if err != errTimeout {
			t.Errorf("%s: error %q is expected, but received %q instead", test.name, errTimeout, err)
		}
		release()
		sem.(semaphore).Flush()
	}
}

func TestSemaphore_Capacity_Immutability(t *testing.T) {
	capacity := 7

	sem := New(capacity)
	defer sem.(semaphore).Flush()

	if sem.Capacity() != capacity {
		t.Errorf("capacity equals to %d is expected, but received %d instead", capacity, sem.Capacity())
	}

	for i := 0; i < sem.Capacity(); i++ {
		_, _ = sem.Acquire(nil)
	}

	if sem.Capacity() != capacity {
		t.Errorf("capacity equals to %d is expected, but received %d instead", capacity, sem.Capacity())
	}
}

func TestSemaphore_Occupied_Linearity(t *testing.T) {
	sem := New(7)
	defer sem.(semaphore).Flush()

	for i := 0; i < sem.Capacity(); i++ {
		if sem.Occupied() != i {
			t.Errorf("%d occupied places are expected, but received %d instead", i, sem.Occupied())
		}
		_, _ = sem.Acquire(nil)
	}

	if sem.Occupied() != sem.Capacity() {
		t.Errorf("%d occupied places are expected, but received %d instead", sem.Capacity(), sem.Occupied())
	}
}

func TestSemaphore_Release_TryToGetDeadLock(t *testing.T) {
	sem := New(0)

	if err := sem.Release(); err != errEmpty {
		t.Errorf("error %q is expected, but received %q instead", errEmpty, err)
	}
}

func TestSemaphore_Concurrently(t *testing.T) {
	sem := New(int(math.Max(2.0, float64(runtime.GOMAXPROCS(0)))))
	defer sem.(semaphore).Flush()

	var counter int32

	start, wg := make(chan bool), &sync.WaitGroup{}
	for i := 0; i < sem.Capacity(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			release, err := sem.Acquire(nil)
			if err != nil {
				t.Errorf("error is not expected, but received %q instead", err)
				return
			}
			defer release()
			atomic.AddInt32(&counter, 1)
		}()
	}
	close(start)
	wg.Wait()

	if int(counter) != sem.Capacity() {
		t.Errorf("counter value equals to %d is expected, but received %d instead", sem.Capacity(), counter)
	}

	if sem.Occupied() != 0 {
		t.Errorf("zero occupied places are expected, but received %d instead", sem.Occupied())
	}
}

func BenchmarkSemaphore_Acquire(b *testing.B) {
	sem := New(b.N)
	defer sem.(semaphore).Flush()

	for i := 0; i < b.N; i++ {
		_, _ = sem.Acquire(nil)
	}

	if sem.Occupied() != sem.Capacity() {
		b.Error("full filled semaphore is expected")
	}
}

func BenchmarkSemaphore_Acquire_Release(b *testing.B) {
	sem := New(b.N)
	defer sem.(semaphore).Flush()

	for i := 0; i < b.N; i++ {
		_, _ = sem.Acquire(nil)
		_ = sem.Release()
	}

	if sem.Occupied() != 0 {
		b.Error("empty semaphore is expected")
	}
}
