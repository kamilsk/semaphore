// +build go1.7

package semaphore_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/kamilsk/semaphore/v4"
)

// This example shows how to use context and semaphore together.
func Example_semaphoreWithContext() {
	deadliner := func(limit int, timeout time.Duration, handler http.HandlerFunc) http.HandlerFunc {
		throughput := New(limit)
		return func(rw http.ResponseWriter, req *http.Request) {
			ctx := WithContext(req.Context(), WithTimeout(timeout))

			release, err := throughput.Acquire(ctx.Done())
			if err != nil {
				http.Error(rw, err.Error(), http.StatusGatewayTimeout)
				return
			}
			defer release()

			handler.ServeHTTP(rw, req.WithContext(ctx))
		}
	}

	ts := httptest.NewServer(deadliner(2, sla, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// context deadline is expected
		select {
		case <-req.Context().Done():
			rw.WriteHeader(http.StatusGatewayTimeout)
		case <-time.After(2 * sla):
			rw.WriteHeader(http.StatusOK)
		}
	})))
	defer ts.Close()

	ok, fail := sendParallelHTTPRequestsToURL(5, ts.URL)

	fmt.Printf("success: %d, failure: %d \n", ok, fail)
	// Output: success: 0, failure: 5
}
