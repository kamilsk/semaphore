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

// Define semaphore for maximum two operations at the same time.
var (
	limit = 2
	sla   = 50 * time.Millisecond
)

// This example shows how to follow SLA.
func ExampleHTTPResponseTimeLimitation() {
	limiter := semaphore.New(limit)

	// start http server to handle parallel requests
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		done := make(chan struct{})

		ctx, cancel := context.WithTimeout(req.Context(), sla)
		defer cancel()

		release, err := limiter.Acquire(ctx)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusGatewayTimeout)
			return
		}
		defer release()

		go func() {
			defer close(done)

			// do some heavy work
			time.Sleep(40 * time.Millisecond)
		}()

		// wait what happens before
		select {
		case <-ctx.Done():
			http.Error(rw, ctx.Err().Error(), http.StatusGatewayTimeout)
			return
		case <-done:
			rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rw.WriteHeader(http.StatusOK)
			return
		}

	}))
	defer ts.Close()

	ok, fail := sendParallelHTTPRequestsToURL(5, ts.URL)

	fmt.Printf("success: %d, failure: %d \n", ok, fail)
	// Output: success: 2, failure: 3
}

// Sends five parallel HTTP requests to the specified URL.
func sendParallelHTTPRequestsToURL(parallelism int, url string) (success, failure int32) {
	start := make(chan bool)
	wg := &sync.WaitGroup{}

	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer func() { _ = resp.Body.Close() }()

			// calculate the result
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt32(&success, 1)
			case http.StatusGatewayTimeout:
				atomic.AddInt32(&failure, 1)
			}
		}()
	}

	close(start)
	wg.Wait()

	return
}
