> # Semaphore
>
> Simple non-blocking semaphore implementation with timeout written on Go.
>
> > Inspired by [go-resiliency](https://github.com/eapache/go-resiliency) and [sema](https://github.com/tarndt/sema).

[![Build Status](https://travis-ci.org/kamilsk/semaphore.svg?branch=master)](https://travis-ci.org/kamilsk/semaphore)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/semaphore/badge.svg)](https://coveralls.io/github/kamilsk/semaphore)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/semaphore)](https://goreportcard.com/report/github.com/kamilsk/semaphore)
[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore?status.svg)](https://godoc.org/github.com/kamilsk/semaphore)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE.md)

## Usages

```go
if err := semaphore.Acquire(50*time.Millisecond); err != nil {
    // try again later
    return
}
defer semaphore.Release()

// push to monitoring semaphore.Occupied()
// do some heavy work
```

## Installation

```bash
$ go get github.com/kamilsk/semaphore
```

### Mirror

```bash
$ go get bitbucket.org/kamilsk/semaphore
```

## Feedback

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kamilsk/semaphore)
[![@ikamilsk](https://img.shields.io/badge/author-%40ikamilsk-blue.svg)](https://twitter.com/ikamilsk)
