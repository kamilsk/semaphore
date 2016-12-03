// +build go1.7

package semaphore_test

import (
	"context"
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
func Example_Limiter() {
	// start http server to handle parallel requests
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), time.Millisecond)
		defer cancel()

		release, err := limiter.Acquire(ctx)
		if err != nil {
			rw.WriteHeader(http.StatusTooManyRequests)
			return
		}
		defer release()

		// do some heavy work
		time.Sleep(time.Second)
		rw.WriteHeader(http.StatusOK)
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
