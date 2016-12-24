package simple_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kamilsk/semaphore/simple"
)

// define semaphore for maximum two operations at the same time
var limiter = simple.New(2)

// This example shows how to restrict an active number of heavy operations.
func Example_limiter() {
	// start http server to handle parallel requests
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if err := limiter.Acquire(50 * time.Millisecond); err != nil {
			rw.WriteHeader(http.StatusTooManyRequests)
			return
		}
		defer func() { _ = limiter.Release() }()

		// do some heavy work
		time.Sleep(250 * time.Millisecond)
		rw.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ok, fail := sendParallelHTTPRequestsToURL(ts.URL)

	fmt.Printf("success: %d, failure: %d \n", ok, fail)
	// Output: success: 2, failure: 3
}

// sends five parallel HTTP requests to the specified url
func sendParallelHTTPRequestsToURL(url string) (success, failure int32) {
	start := make(chan bool)
	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// calculate the result
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt32(&success, 1)
			case http.StatusTooManyRequests:
				atomic.AddInt32(&failure, 1)
			}
		}()
	}

	close(start)
	wg.Wait()

	return
}
