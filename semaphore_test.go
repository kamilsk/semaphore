package semaphore_test

import (
	"context"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/kamilsk/semaphore/v5"
	"github.com/stretchr/testify/assert"
)

func TestSemaphore_Acquire_Timeout(t *testing.T) {
	for _, tc := range []struct {
		name    string
		timeout time.Duration
	}{
		{name: "negative timeout", timeout: -time.Second},
		{name: "zero timeout", timeout: 0},
		{name: "positive timeout", timeout: time.Nanosecond},
	} {
		semaphore := New(0)
		ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
		release, err := semaphore.Acquire(ctx.Done())
		cancel()
		_ = release.Release()

		assert.True(t, IsTimeout(err))
	}
}

func TestSemaphore_Capacity_Immutability(t *testing.T) {
	capacity := 7

	semaphore := New(capacity)

	if semaphore.Capacity() != capacity {
		t.Errorf("an unexpected capacity. expected: %d; obtained: %d", capacity, semaphore.Capacity())
	}

	for i := 0; i < semaphore.Capacity(); i++ {
		_, _ = semaphore.Acquire(nil)
	}

	if semaphore.Capacity() != capacity {
		t.Errorf("an unexpected capacity. expected: %d; obtained: %d", capacity, semaphore.Capacity())
	}
}

func TestSemaphore_Occupied_Linearity(t *testing.T) {
	semaphore := New(7)

	for i := 0; i < semaphore.Capacity(); i++ {
		if semaphore.Occupied() != i {
			t.Errorf("unexpected occupied places. expected: %d; obtained: %d", i, semaphore.Occupied())
		}
		_, _ = semaphore.Acquire(nil)
	}

	if semaphore.Occupied() != semaphore.Capacity() {
		t.Errorf("unexpected occupied places. expected: %d; obtained: %d", semaphore.Capacity(), semaphore.Occupied())
	}
}

func TestSemaphore_Release_TryToGetDeadLock(t *testing.T) {
	semaphore := New(0)

	if err, expected := semaphore.Release(), "semaphore is empty"; !IsEmpty(err) {
		t.Errorf("an unexpected error. expected: %s; obtained: %v", expected, err)
	}
}

func TestSemaphore_Signal(t *testing.T) {
	semaphore := New(0)

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	release, ok := <-semaphore.Signal(ctx.Done())
	cancel()
	if release != nil || ok {
		t.Error("unexpected signal")
	}
}

func TestSemaphore_Concurrently(t *testing.T) {
	semaphore := New(int(math.Max(2.0, float64(runtime.GOMAXPROCS(0)))))

	var counter int32

	start, wg := make(chan bool), &sync.WaitGroup{}
	for i := 0; i < semaphore.Capacity(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			release, err := semaphore.Acquire(nil)
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

	if int(counter) != semaphore.Capacity() {
		t.Errorf("an unexpected counter value. expected: %d; obtained: %d", semaphore.Capacity(), counter)
	}

	if semaphore.Occupied() != 0 {
		t.Errorf("zero occupied places are expected but received %d instead", semaphore.Occupied())
	}
}
