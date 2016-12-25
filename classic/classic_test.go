package classic

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
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

	var v int
	from, to := []int{1, 2, 3}, make([]int, 0, 3)
	expected := []int{1, 2, 3}

	wg := sync.WaitGroup{}
	for i := 0; i < cap(expected); i++ {
		wg.Add(1)
		go func() {
			sem.Lock()
			v, from = from[0], from[1:]
			to = append(to, v)
			sem.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	if !reflect.DeepEqual(expected, to) {
		t.Errorf("%+v is not equal to %+v", to, expected)
	}
}
