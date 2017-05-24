package semaphore_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/kamilsk/semaphore"
)

// This example shows how to limit request' throughput.
func ExampleHTTPRequestThroughputLimitation() {
	limiter := func(limit int, timeout time.Duration, handler http.Handler) http.Handler {
		throughput := semaphore.New(limit)
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(req.Context(), timeout)
			defer cancel()

			release, err := throughput.Acquire(ctx)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusTooManyRequests)
				return
			}
			defer release()

			handler.ServeHTTP(rw, req)
		})
	}

	var race int
	ts := httptest.NewServer(limiter(1, sla, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// do some limited work
		race++
	})))
	defer ts.Close()

	ok, fail := sendParallelHTTPRequestsToURL(5, ts.URL)

	fmt.Printf("success: %d, failure: %d, race: %t \n", ok, fail, race != 5)
	// Output: success: 5, failure: 0, race: false
}
