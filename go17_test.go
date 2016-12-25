// +build go1.7

package semaphore

import (
	"context"
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

func TestSemaphore_Acquire_InvalidTimeout(t *testing.T) {
	sem := New(0)
	defer sem.(semaphore).Flush()

	nothingToDo := func(context.CancelFunc) {}

	for _, test := range []struct {
		name    string
		timeout time.Duration
		do      func(cancel context.CancelFunc)
	}{
		{name: "negative timeout", timeout: -time.Second, do: nothingToDo},
		{name: "zero timeout", timeout: 0, do: nothingToDo},
		{name: "positive timeout", timeout: time.Nanosecond, do: nothingToDo},
		{name: "context cancel", timeout: time.Second, do: func(cancel context.CancelFunc) { cancel() }},
	} {
		ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
		test.do(cancel)
		release, err := sem.Acquire(ctx)
		if err != errTimeout {
			t.Errorf("%s: error %q is expected, but received %q instead", test.name, errTimeout, err)
		}
		release()
		cancel()
	}
}

func TestSemaphore_Capacity_Immutability(t *testing.T) {
	capacity := 7

	sem := New(capacity)
	defer sem.(semaphore).Flush()

	if sem.Capacity() != capacity {
		t.Errorf("capacity equal to %d is expected, but received %d instead", capacity, sem.Capacity())
	}

	ctx := context.Background()
	for i := 0; i < sem.Capacity(); i++ {
		_, _ = sem.Acquire(ctx)
	}

	if sem.Capacity() != capacity {
		t.Errorf("capacity equal to %d is expected, but received %d instead", capacity, sem.Capacity())
	}
}

func TestSemaphore_Occupied_Linearity(t *testing.T) {
	sem := New(7)
	defer sem.(semaphore).Flush()

	ctx := context.Background()
	for i := 0; i < sem.Capacity(); i++ {
		if sem.Occupied() != i {
			t.Errorf("%d occupied places are expected, but received %d instead", i, sem.Occupied())
		}
		_, _ = sem.Acquire(ctx)
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
	for i, ctx := 0, context.Background(); i < sem.Capacity(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			release, err := sem.Acquire(ctx)
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
		t.Errorf("counter value equal to %d is expected, but received %d instead", sem.Capacity(), counter)
	}

	if sem.Occupied() != 0 {
		t.Errorf("zero occupied places are expected, but received %d instead", sem.Occupied())
	}
}

func BenchmarkSemaphore_Acquire(b *testing.B) {
	ctx, sem := context.Background(), New(b.N)
	defer sem.(semaphore).Flush()

	for i := 0; i < b.N; i++ {
		_, _ = sem.Acquire(ctx)
	}

	if sem.Occupied() != sem.Capacity() {
		b.Error("full filled semaphore is expected")
	}
}

func BenchmarkSemaphore_Acquire_Release(b *testing.B) {
	ctx, sem := context.Background(), New(b.N)
	defer sem.(semaphore).Flush()

	for i := 0; i < b.N; i++ {
		_, _ = sem.Acquire(ctx)
		_ = sem.Release()
	}

	if sem.Occupied() != 0 {
		b.Error("empty semaphore is expected")
	}
}
