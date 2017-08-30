> # semaphore
>
> Semaphore pattern implementation with a timeout of lock/unlock operations based on channels.

[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#goroutines)
[![Build Status](https://travis-ci.org/kamilsk/semaphore.svg?branch=master)](https://travis-ci.org/kamilsk/semaphore)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/semaphore/badge.svg)](https://coveralls.io/github/kamilsk/semaphore)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/semaphore)](https://goreportcard.com/report/github.com/kamilsk/semaphore)
[![Exago](https://api.exago.io/badge/rank/github.com/kamilsk/semaphore)](https://www.exago.io/project/github.com/kamilsk/semaphore)
[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore?status.svg)](https://godoc.org/github.com/kamilsk/semaphore)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE)

## Code of Conduct

The project team follows [Contributor Covenant v1.4](http://contributor-covenant.org/version/1/4/).
Instances of abusive, harassing or otherwise unacceptable behavior may be reported by contacting
the project team at feedback@octolab.org.

---

## Usage

### Console tool for command execution in parallel

This example shows how to execute many console commands in parallel.

```bash
$ semaphore create 2
$ semaphore add -- docker build
$ semaphore add -- vagrant up
$ semaphore add -- ansible-playbook
$ semaphore wait --notify --timeout=1m
```

See more details [here](cmd#semaphore).

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

See more details [here](https://godoc.org/github.com/kamilsk/semaphore#example-package--HttpResponseTimeLimitation).

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

See more details [here](https://godoc.org/github.com/kamilsk/semaphore#example-package--HttpRequestThroughputLimitation).

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

See more details [here](https://godoc.org/github.com/kamilsk/semaphore#example-package--SemaphoreWithContext).

### Pool of workers

```go
type Pool struct {
	sem  semaphore.Semaphore
	work chan func()
}

func (p *Pool) Schedule(task func()) {
	select {
	case p.work <- task:
	case <-p.sem.Signal(nil):
		go p.worker(task)
	}
}

func (p *Pool) worker(task func()) {
	defer func() { p.sem.Release() }()
	var ok bool
	for {
		task()
		task, ok = <-p.work
		if !ok {
			return
		}
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
	pool.Schedule(func() { fmt.Println(1) })
	pool.Schedule(func() { fmt.Println(2) })
	pool.Schedule(func() { fmt.Println(3) })
	pool.Schedule(func() { fmt.Println(4) })
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
```

### Mirror

```bash
$ egg bitbucket.org/kamilsk/semaphore
```

> [egg](https://github.com/kamilsk/egg) is an `extended go get`.

### Update

This library is using [SemVer](http://semver.org) for versioning and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe.
Therefore, do not use `go get -u` to update it, use [Glide](https://glide.sh) or something similar for this purpose.

## Contributing workflow

Read first [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Code quality checking

```bash
$ make docker-pull-tools
$ make check-code-quality
```

### Testing

#### Local

```bash
$ make install-deps
$ make test # or test-with-coverage
$ make bench
```

#### Docker

```bash
$ make docker-pull
$ make complex-tests # or complex-tests-with-coverage
$ make complex-bench
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/semaphore)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)

## Notes

- tested on Go 1.5, 1.6, 1.7, 1.8 and 1.9
- [research](research)
