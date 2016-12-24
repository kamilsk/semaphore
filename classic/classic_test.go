package classic

import (
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func (sem semaphore) Flush() {
	close(sem)
	for range sem {
	}
}

func TestLockingSemaphore(t *testing.T) {
	data := []int32{1, 2, 3}

	sem := NewLocking(len(data))
	defer sem.(semaphore).Flush()

	var sum int32

	for _, i := range data {
		go func(delta int32) {
			atomic.AddInt32(&sum, delta)
			sem.P(1)
		}(i)
	}
	sem.V(len(data))

	if int(sum) != 6 {
		t.Errorf("sum equal to 6 is expected, but received %d instead", sum)
	}
}

func TestSyncingSemaphore(t *testing.T) {
	data := []int32{1, 2, 3}

	sem := NewSyncing(len(data))
	defer sem.(semaphore).Flush()

	var sum int32

	for _, i := range data {
		go func(delta int32) {
			atomic.AddInt32(&sum, delta)
			sem.Signal()
		}(i)
	}
	sem.Wait(len(data))

	if int(sum) != 6 {
		t.Errorf("sum equal to 6 is expected, but received %d instead", sum)
	}
}

func TestBinarySemaphore(t *testing.T) {
	sem := NewBinary()
	defer sem.(semaphore).Flush()

	expected, steps := []string{"first", "second", "third"}, make([]string, 0, 3)

	sem.Lock()
	go func() {
		defer sem.Unlock()
		sem.Lock()

		steps = append(steps, "second")

		go func() {
			defer sem.Unlock()
			sem.Lock()

			steps = append(steps, "third")
		}()
	}()
	steps = append(steps, "first")
	sem.Unlock()

	// just enough to yield the scheduler and let the goroutines work off
	time.Sleep(time.Millisecond)

	sem.Lock()
	if !reflect.DeepEqual(expected, steps) {
		t.Errorf("%+v is not equal to %+v", steps, expected)
	}
	sem.Unlock()
}
