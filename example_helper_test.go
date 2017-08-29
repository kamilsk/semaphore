package semaphore_test

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const sla = 50 * time.Millisecond

// Sends parallel HTTP requests to the specified URL.
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
