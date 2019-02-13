// +build go1.6

package semaphore_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/kamilsk/semaphore/v4"
)

// This example shows how to limit request throughput.
func Example_httpRequestThroughputLimitation() {
	limiter := func(limit int, timeout time.Duration, handler http.HandlerFunc) http.HandlerFunc {
		throughput := New(limit)
		return func(rw http.ResponseWriter, req *http.Request) {
			deadline := WithTimeout(timeout)

			release, err := throughput.Acquire(deadline)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusTooManyRequests)
				return
			}
			defer release()

			handler.ServeHTTP(rw, req)
		}
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
