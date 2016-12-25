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
var (
	limiter   = semaphore.New(2)
	sla       = 50 * time.Millisecond
	timeIsOut = func(rw http.ResponseWriter, err error) {
		http.Error(rw, err.Error(), http.StatusGatewayTimeout)
	}
)

// This example shows how to follow SLA.
func Example_sla() {
	// start http server to handle parallel requests
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		done := make(chan struct{})

		// user defined timeout
		timeout, err := time.ParseDuration(req.FormValue("timeout"))
		if err != nil || sla < timeout {
			timeout = sla
		}

		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		defer cancel()

		release, err := limiter.Acquire(ctx)
		if err != nil {
			timeIsOut(rw, err)
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
			timeIsOut(rw, ctx.Err())
			return
		case <-done:
			rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rw.WriteHeader(http.StatusOK)
			return
		}

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
			case http.StatusGatewayTimeout:
				atomic.AddInt32(&failure, 1)
			}
		}()
	}

	close(start)
	wg.Wait()

	return
}
