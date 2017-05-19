> # semaphore
>
> Semaphore pattern implementation with timeout of lock/unlock operations based on channel and context.

[![Build Status](https://travis-ci.org/kamilsk/semaphore.svg?branch=master)](https://travis-ci.org/kamilsk/semaphore)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/semaphore/badge.svg)](https://coveralls.io/github/kamilsk/semaphore)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/semaphore)](https://goreportcard.com/report/github.com/kamilsk/semaphore)
[![Exago](https://api.exago.io/badge/rank/github.com/kamilsk/semaphore)](https://www.exago.io/project/github.com/kamilsk/semaphore)
[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore?status.svg)](https://godoc.org/github.com/kamilsk/semaphore)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE)

## Usage

### HTTP response' time limitation

```go
sla := 100 * time.Millisecond
sem := semaphore.New(1000)

http.Handle("/do-with-timeout", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
    done := make(chan struct{})

    ctx, cancel := context.WithTimeout(req.Context(), sla)
    defer cancel()

    release, err := sem.Acquire(ctx)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusGatewayTimeout)
        return
    }
    defer release()

    go func() {
        defer close(done)

        // do some heavy work
    }()

    // wait what happens before
    select {
    case <-ctx.Done():
        http.Error(rw, err.Error(), http.StatusGatewayTimeout)
    case <-done:
        // send success response
    }
}))
```

### HTTP request' throughput limitation

```go
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

		handler(rw, req)
	})
}

http.Handle("/do-limited", limiter(1000, time.Minute, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	// do some limited work
})))
```

## Installation

```bash
$ go get github.com/kamilsk/semaphore
```

### Mirror

```bash
$ egg -fix-vanity-url bitbucket.org/kamilsk/semaphore
```

> [egg](https://github.com/kamilsk/egg) is an `extended go get`.

### Update

This library is using [SemVer](http://semver.org) for versioning and it is not
[BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe.
Therefore, do not use `go get -u` to update it, use [Glide](https://glide.sh) or something similar for this purpose.

## Contributing workflow

### Code quality checking

```bash
$ make docker-pull-tools
$ make docker-gometalinter
```

### Testing

#### Local

```bash
$ make install-deps
$ make test-with-coverage
```

#### Docker

```bash
$ make docker-pull
$ make docker-test-with-coverage
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/semaphore)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)

## Notes

- tested on Go 1.7 and 1.8, use 2.x version for 1.5 and 1.6
- [research](RESEARCH.md)
