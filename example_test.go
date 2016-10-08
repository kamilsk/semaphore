package semaphore

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"
)

// This example shows how to restrict the active number of heavy operations.
func Example() {
	// maximum two operations the same time
	limiter := New(2)

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if err := limiter.Acquire(time.Millisecond); err != nil {
			rw.WriteHeader(429) // http.StatusTooManyRequests
			_, _ = rw.Write([]byte("please try again later"))
			return
		}
		defer limiter.Release()

		// do some heavy work
		time.Sleep(time.Second)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("success"))
	}))
	defer ts.Close()

	start := make(chan bool)
	wg := &sync.WaitGroup{}
	var ok, fail int32
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			resp, err := http.Get(ts.URL)
			if err != nil {
				return
			}
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt32(&ok, 1)
			case 429: // http.StatusTooManyRequests
				atomic.AddInt32(&fail, 1)
			}
		}()
	}
	close(start)
	wg.Wait()

	// as expected "success" shows 2, "failure" shows 3
	fmt.Printf("success: %d, failure: %d\n", ok, fail)
}
