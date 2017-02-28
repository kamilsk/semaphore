> # semaphore
>
> Semaphore pattern implementation with timeout of lock/unlock operations based on channel and context.

[![Build Status](https://travis-ci.org/kamilsk/semaphore.svg?branch=master)](https://travis-ci.org/kamilsk/semaphore)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/semaphore/badge.svg)](https://coveralls.io/github/kamilsk/semaphore)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/semaphore)](https://goreportcard.com/report/github.com/kamilsk/semaphore)
[![Exago](https://api.exago.io/badge/rank/github.com/kamilsk/semaphore)](https://exago.io/project/github.com/kamilsk/semaphore)
[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore?status.svg)](https://godoc.org/github.com/kamilsk/semaphore)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE)

## Usage

### HTTP request limitation with timeout response

```go
sla := 100 * time.Millisecond
sem := semaphore.New(1000)

timeIsOver := func(rw http.ResponseWriter, err error) {
    http.Error(rw, err.Error(), http.StatusGatewayTimeout)
}

http.Handle("/do-with-timeout", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
    done := make(chan struct{})

    // user defined timeout
    timeout, err := time.ParseDuration(req.FormValue("timeout"))
    if err != nil || sla < timeout {
        timeout = sla
    }

    ctx, cancel := context.WithTimeout(req.Context(), timeout)
    defer cancel()

    release, err := sem.Acquire(ctx)
    if err != nil {
        timeIsOver(rw, err)
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
        timeIsOver(rw, ctx.Err())
    case <-done:
        rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
        rw.WriteHeader(http.StatusOK)
    }
}))
```

## Installation

```bash
$ egg -version 2.x github.com/kamilsk/semaphore
```

### Mirror

```bash
$ go get bitbucket.org/kamilsk/semaphore | egg -fix-vanity-url -version 2.x
```

> [egg](https://github.com/kamilsk/egg) is an `extended go get`.

### Update

This library is using [SemVer](http://semver.org) for versioning and it is not [BC](https://en.wikipedia.org/wiki/Backward_compatibility)-safe.
Therefore, do not use `go get -u` to update it, use [Glide](https://glide.sh) or something similar for this purpose.

## Integration with Docker

```bash
$ make docker-pull
$ make docker-gometalinter ARGS=--deadline=12s
$ make docker-bench ARGS=-benchmem
$ make docker-test ARGS=-v
$ make docker-test-with-coverage ARGS=-v OPEN_BROWSER=true
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/semaphore)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)

## Notes

- tested on Go 1.5, 1.6, 1.7 and 1.8
- [research](RESEARCH.md)
