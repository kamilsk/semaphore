package classic

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestBlockingSemaphore(t *testing.T) {
	data := []int32{0, 1, 0, 1}
	sem := NewBlockingSemaphore(len(data))
	var sum int32

	for _, i := range data {
		go func(delta int32) {
			atomic.AddInt32(&sum, delta)
			sem.P(1)
		}(i)
	}
	sem.V(len(data))

	if sum != int32(2) {
		t.Errorf("expected sum value is equals to 2, obtained %d", sum)
	}
}

func TestProcessSemaphore(t *testing.T) {
	data := []int32{0, 1, 0, 1}
	sem := NewProcessSemaphore(len(data))
	var sum int32

	for _, i := range data {
		go func(delta int32) {
			atomic.AddInt32(&sum, delta)
			sem.Signal()
		}(i)
	}
	sem.Wait(len(data))

	if sum != int32(2) {
		t.Errorf("expected sum value is equals to 2, obtained %d", sum)
	}
}

func TestBinarySemaphore(t *testing.T) {
	sem := NewBinarySemaphore()
	var step int

	go func() {
		defer sem.Unlock()
		sem.Lock()
		if step != 1 {
			t.Fatal("unexpected result")
		}
		step = 2
	}()

	sem.Lock()
	if step != 0 {
		t.Fatal("unexpected result")
	}
	step = 1
	sem.Unlock()

	// just enough to yield the scheduler and let the goroutines work off
	time.Sleep(time.Millisecond)

	sem.Lock()
	if step != 2 {
		t.Fatal("unexpected result")
	}
}
