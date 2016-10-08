package semaphore

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSemaphore_Concurrently(t *testing.T) {
	size := int(math.Max(2.0, float64(runtime.GOMAXPROCS(0))))
	sem := New(size)
	var counter int32

	start := make(chan bool)
	wg := &sync.WaitGroup{}
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if err := sem.Acquire(time.Millisecond); err != nil {
				t.Fatal("error is not expected")
				return
			}
			defer sem.Release()
			atomic.AddInt32(&counter, 1)
		}()
	}
	close(start)
	wg.Wait()

	if counter != int32(size) {
		t.Errorf("expected counter value is equals to %d, obtained %d", size, counter)
	}
}

func TestTimeoutError_Concurrently(t *testing.T) {
	size := int(math.Max(2.0, float64(runtime.GOMAXPROCS(0))))
	sem := New(size)
	var counter int32

	start := make(chan bool)
	wg := &sync.WaitGroup{}
	for i := 0; i < size+1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if err := sem.Acquire(time.Millisecond); err != nil {
				atomic.AddInt32(&counter, 1)
				if err.Error() != errTimeout.Error() {
					t.Errorf("expected error message %q, obtained %q", err, errTimeout)
				}
				return
			}
			defer sem.Release()
			time.Sleep(time.Second)
		}()
	}
	close(start)
	wg.Wait()

	if counter != 1 {
		t.Errorf("expected counter value is equals to %d, obtained %d", 1, counter)
	}
}
