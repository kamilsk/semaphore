// +build go1.7

package semaphore_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"sync"
	"time"

	. "github.com/kamilsk/semaphore/v4"
)

// User represents user ID.
type User int

// Config contains abstract configuration fields.
type Config struct {
	DefaultCapacity  int
	DefaultRateLimit time.Duration
	DefaultUser      User
	Capacity         map[User]int
	RateLimit        map[User]time.Duration
	SLA              time.Duration
}

// This variables can be a part of limiter provider service.
var (
	mx       sync.RWMutex
	limiters = make(map[User]Semaphore)
)

// LimiterForUser returns limiter for user found in request context.
func LimiterForUser(user User, cnf Config) Semaphore {
	mx.RLock()
	limiter, ok := limiters[user]
	mx.RUnlock()
	if !ok {
		mx.Lock()
		limiter, ok = limiters[user]
		if !ok {
			c, ok := cnf.Capacity[user]
			if !ok {
				c = cnf.DefaultCapacity
			}
			limiter = New(c)
			limiters[user] = limiter
		}
		mx.Unlock()
	}
	return limiter
}

// RateLimiter performs rate limitation.
func RateLimiter(cnf Config, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		user, ok := req.Context().Value("user").(User)
		if !ok {
			user = cnf.DefaultUser
		}

		limiter := LimiterForUser(user, cnf)
		release, err := limiter.Acquire(WithTimeout(cnf.SLA))
		if err != nil {
			http.Error(rw, "operation timeout", http.StatusGatewayTimeout)
			return
		}

		go func() { handler.ServeHTTP(rw, req) }()

		rl, ok := cnf.RateLimit[user]
		if !ok {
			rl = cnf.DefaultRateLimit
		}
		time.Sleep(rl)
		release()
	}
}

// UserToContext gets user ID from request header and puts it into request context.
func UserToContext(cnf Config, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var user User = cnf.DefaultUser

		if id := req.Header.Get("user"); id != "" {
			i, err := strconv.Atoi(id)
			if err == nil {
				user = User(i)
			}
		}

		handler.ServeHTTP(rw, req.WithContext(context.WithValue(req.Context(), "user", user)))
	}
}

// This example shows how to create user-specific rate limiter.
func Example_userRateLimitation() {
	var cnf Config = Config{
		DefaultUser:      1,
		DefaultCapacity:  runtime.GOMAXPROCS(0),
		DefaultRateLimit: 10 * time.Millisecond,
		Capacity:         map[User]int{1: 1},
		RateLimit:        map[User]time.Duration{1: 100 * time.Millisecond},
		SLA:              time.Second,
	}

	ts := httptest.NewServer(RateLimiter(cnf, UserToContext(cnf, func(rw http.ResponseWriter, req *http.Request) {})))
	defer ts.Close()

	start := time.Now()
	ok, fail := sendParallelHTTPRequestsToURL(2, ts.URL)
	end := time.Now()

	fmt.Printf("success: %d, failure: %d, elapsed: ~%d ms \n", ok, fail, (end.Sub(start).Nanoseconds()/100000000)*100)
	// Output: success: 2, failure: 0, elapsed: ~200 ms
}
