package semaphore_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kamilsk/semaphore"
)

// This example shows how to follow SLA.
func Example_httpResponseTimeLimitation() {
	limiter := semaphore.New(2)

	// start http server to handle parallel requests
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		done := make(chan struct{})
		deadline := semaphore.WithTimeout(sla)

		go func() {
			release, err := limiter.Acquire(deadline)
			if err != nil {
				return
			}
			defer release()
			defer close(done)

			// do some heavy work
			time.Sleep(40 * time.Millisecond)
		}()

		// wait what happens before
		select {
		case <-deadline:
			http.Error(rw, "operation timeout", http.StatusGatewayTimeout)
		case <-done:
			// send success response
			rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rw.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()

	ok, fail := sendParallelHTTPRequestsToURL(5, ts.URL)

	fmt.Printf("success: %d, failure: %d \n", ok, fail)
	// Output: success: 2, failure: 3
}
