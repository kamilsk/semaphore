> # 🚦 semaphore [![Tweet][icon_twitter]][twitter_publish]
> [![Analytics][analytics_pixel]][page_promo]
> Semaphore pattern implementation with timeout of lock/unlock operations based on channels.

[![Awesome][icon_awesome]](https://github.com/avelino/awesome-go#goroutines)
[![Patreon][icon_patreon]](https://www.patreon.com/octolab)
[![Build Status][icon_build]][page_build]
[![Code Coverage][icon_coverage]][page_quality]
[![Code Quality][icon_quality]][page_quality]
[![GoDoc][icon_docs]][page_docs]
[![Research][icon_research]](../../tree/research)
[![License][icon_license]](LICENSE)

## Usage

### Quick start

```go
limiter := semaphore.New(1000)

http.HandleFunc("/", func(rw http.ResponseWriter, _ *http.Request) {
	if _, err := limiter.Acquire(semaphore.WithTimeout(time.Minute)); err != nil {
		http.Error(rw, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	defer limiter.Release()
	// handle request
})

log.Fatal(http.ListenAndServe(":80", nil))
```

### Console tool for command execution in parallel

This example shows how to execute many console commands in parallel.

```bash
$ semaphore create 2
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait --timeout=1m --notify
```

[![asciicast](https://asciinema.org/a/136111.png)](https://asciinema.org/a/136111)

See more details [here](cmd/semaphore).

### HTTP response time limitation

This example shows how to follow SLA.

```go
sla := 100 * time.Millisecond
sem := semaphore.New(1000)

http.Handle("/do-with-timeout", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	done := make(chan struct{})
	deadline := semaphore.WithTimeout(sla)

	go func() {
		release, err := sem.Acquire(deadline)
		if err != nil {
			return
		}
		defer release()
		defer close(done)

		// do some heavy work
	}()

	// wait what happens before
	select {
	case <-deadline:
		http.Error(rw, "operation timeout", http.StatusGatewayTimeout)
	case <-done:
		// send success response
	}
}))
```

### HTTP request throughput limitation

This example shows how to limit request throughput.

```go
limiter := func(limit int, timeout time.Duration, handler http.HandlerFunc) http.HandlerFunc {
	throughput := semaphore.New(limit)
	return func(rw http.ResponseWriter, req *http.Request) {
		deadline := semaphore.WithTimeout(timeout)

		release, err := throughput.Acquire(deadline)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusTooManyRequests)
			return
		}
		defer release()

		handler.ServeHTTP(rw, req)
	}
}

http.HandleFunc("/do-with-limit", limiter(1000, time.Minute, func(rw http.ResponseWriter, req *http.Request) {
	// do some limited work
}))
```

### HTTP personal rate limitation

This example shows how to create user-specific rate limiter.

```go
func LimiterForUser(user User, cnf Config) semaphore.Semaphore {
	mx.RLock()
	limiter, ok := limiters[user]
	mx.RUnlock()
	if !ok {
		mx.Lock()
		// handle negative case
		mx.Unlock()
	}
	return limiter
}

func RateLimiter(cnf Config, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		user, ok := // get user from request context

		limiter := LimiterForUser(user, cnf)
		release, err := limiter.Acquire(semaphore.WithTimeout(cnf.SLA))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusGatewayTimeout)
			return
		}

		// handle the request in separated goroutine because the current will be held
		go func() { handler.ServeHTTP(rw, req) }()

		// hold the place for a required time
		rl, ok := cnf.RateLimit[user]
		if !ok {
			rl = cnf.DefaultRateLimit
		}
		time.Sleep(rl)
		release()
		// rate limit = semaphore capacity / rate limit time, e.g. 10 request per second 
	}
}

http.HandleFunc("/do-with-rate-limit", RateLimiter(cnf, func(rw http.ResponseWriter, req *http.Request) {
	// do some rate limited work
}))
```

### Use context for cancellation

This example shows how to use context and semaphore together.

```go
deadliner := func(limit int, timeout time.Duration, handler http.HandlerFunc) http.HandlerFunc {
	throughput := semaphore.New(limit)
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := semaphore.WithContext(req.Context(), semaphore.WithTimeout(timeout))

		release, err := throughput.Acquire(ctx.Done())
		if err != nil {
			http.Error(rw, err.Error(), http.StatusGatewayTimeout)
			return
		}
		defer release()

		handler.ServeHTTP(rw, req.WithContext(ctx))
	}
}

http.HandleFunc("/do-with-deadline", deadliner(1000, time.Minute, func(rw http.ResponseWriter, req *http.Request) {
	// do some limited work
}))
```

### A pool of workers

This example shows how to create a pool of workers based on semaphore.

```go
type Pool struct {
	sem  semaphore.Semaphore
	work chan func()
}

func (p *Pool) Schedule(task func()) {
	select {
	case p.work <- task: // delay the task to already running workers
	case release, ok := <-p.sem.Signal(nil): if ok { go p.worker(task, release) } // ok is always true in this case
	}
}

func (p *Pool) worker(task func(), release semaphore.ReleaseFunc) {
	defer release()
	var ok bool
	for {
		task()
		task, ok = <-p.work
		if !ok { return }
	}
}

func New(size int) *Pool {
	return &Pool{
		sem:  semaphore.New(size),
		work: make(chan func()),
	}
}

func main() {
	pool := New(2)
	pool.Schedule(func() { /* do some work */ })
	...
	pool.Schedule(func() { /* do some work */ })
}
```

### Interrupt execution

```go
sem := semaphore.New(runtime.GOMAXPROCS(0))
interrupter := semaphore.Multiplex(
	semaphore.WithTimeout(time.Second),
	semaphore.WithSignal(os.Interrupt),
)
_, err := sem.Acquire(interrupter)
if err == nil {
	panic("press Ctrl+C")
}
// successful interruption
```

## Installation

```bash
$ go get github.com/kamilsk/semaphore
$ # or use mirror
$ egg bitbucket.org/kamilsk/semaphore
```

> [egg](https://github.com/kamilsk/egg)<sup id="anchor-egg">[1](#egg)</sup> is an `extended go get`.

## Update

This library is using [SemVer](http://semver.org) for versioning, and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe. Therefore, do not use `go get -u` to update it,
use [dep](https://github.com/golang/dep) or something similar for this purpose.

<sup id="egg">1</sup> The project is still in prototyping. [↩](#anchor-egg)

---

[![Gitter][icon_gitter]](https://gitter.im/kamilsk/semaphore)
[![@kamilsk][icon_tw_author]](https://twitter.com/ikamilsk)
[![@octolab][icon_tw_sponsor]](https://twitter.com/octolab_inc)

made with ❤️ by [OctoLab](https://www.octolab.org/)

[analytics_pixel]: https://ga-beacon.appspot.com/UA-109817251-2/semaphore/dev?pixel

[icon_awesome]:    https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg
[icon_build]:      https://travis-ci.org/kamilsk/semaphore.svg?branch=dev
[icon_coverage]:   https://scrutinizer-ci.com/g/kamilsk/semaphore/badges/coverage.png?b=dev
[icon_docs]:       https://godoc.org/github.com/kamilsk/semaphore?status.svg
[icon_gitter]:     https://badges.gitter.im/Join%20Chat.svg
[icon_license]:    https://img.shields.io/badge/license-MIT-blue.svg
[icon_patreon]:    https://img.shields.io/badge/patreon-donate-orange.svg
[icon_quality]:    https://scrutinizer-ci.com/g/kamilsk/semaphore/badges/quality-score.png?b=dev
[icon_research]:   https://img.shields.io/badge/research-in%20progress-yellow.svg
[icon_tw_author]:  https://img.shields.io/badge/author-%40kamilsk-blue.svg
[icon_tw_sponsor]: https://img.shields.io/badge/sponsor-%40octolab-blue.svg
[icon_twitter]:    https://img.shields.io/twitter/url/http/shields.io.svg?style=social

[page_build]:      https://travis-ci.org/kamilsk/semaphore
[page_docs]:       https://godoc.org/github.com/kamilsk/semaphore
[page_promo]:      https://github.com/kamilsk/semaphore
[page_quality]:    https://scrutinizer-ci.com/g/kamilsk/semaphore/?branch=dev

[twitter_publish]: https://twitter.com/intent/tweet?text=Semaphore%20pattern%20implementation%20with%20a%20timeout%20of%20lock%2Funlock%20operations%20based%20on%20channels&url=https://github.com/kamilsk/semaphore&via=ikamilsk&hashtags=go,semaphore,throughput,limiter
