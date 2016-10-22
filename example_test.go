package semaphore_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kamilsk/semaphore"
)

// define semaphore for maximum two operations at the same time
var limiter = semaphore.New(2)

// This example shows how to restrict the active number of heavy operations.
func Example() {
	// start http server to handle parallel requests
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if err := limiter.Acquire(time.Millisecond); err != nil {
			rw.WriteHeader(429) // http.StatusTooManyRequests
			_, _ = rw.Write([]byte("please try again later"))
			return
		}
		defer func() { _ = limiter.Release() }()

		// do some heavy work
		time.Sleep(time.Second)
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("success"))
	}))
	defer ts.Close()

	ok, fail := sendParallelHTTPRequestsToURL(ts.URL)

	fmt.Printf("success: %d, failure: %d\n", ok, fail)
	// Output: success: 2, failure: 3
}

// sends five HTTP requests at the specified url in parallel
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

			// calculate the result
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt32(&success, 1)
			case 429: // http.StatusTooManyRequests
				atomic.AddInt32(&failure, 1)
			}
		}()
	}

	close(start)
	wg.Wait()

	return
}
