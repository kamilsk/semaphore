> # Semaphore
>
> Simple non-blocking semaphore implementation with timeout written on Go.
>
> > Inspired by
> > - [go-resiliency/semaphore.Semaphore](https://github.com/eapache/go-resiliency/blob/008c74ab45c7c8efbbf0370fdadcf3564faa1e3e/semaphore/semaphore.go)
> > - [godropbox/sync2.Semaphore](https://github.com/dropbox/godropbox/blob/ece4db8e7759e0231f401202ffda6b5946a37ac0/sync2/semaphore.go)
> > - [sema.TimeoutCountingSema](https://github.com/tarndt/sema/blob/02de9df47f0b98e4529584d0b52baa37c2c86e7a/sema.go)
> >
> > Related
> > [Semaphore Pattern](http://tmrts.com/go-patterns/synchronization/semaphore.html)
> > [Semaphore (programming)](https://en.wikipedia.org/wiki/Semaphore_(programming))

[![Build Status](https://travis-ci.org/kamilsk/semaphore.svg?branch=master)](https://travis-ci.org/kamilsk/semaphore)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/semaphore/badge.svg)](https://coveralls.io/github/kamilsk/semaphore)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/semaphore)](https://goreportcard.com/report/github.com/kamilsk/semaphore)
[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore?status.svg)](https://godoc.org/github.com/kamilsk/semaphore)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE.md)

## Usages

```go
sem := semaphore.New(5)

if err := sem.Acquire(50*time.Millisecond); err != nil {
    // try again later
    return
}
defer sem.Release()

// do some heavy work
```

## Tips and tricks

### Retry to acquire a few times with line breaks between attempts

```go
import (
	"github.com/kamilsk/retry"
	"github.com/kamilsk/retry/backoff"
	"github.com/kamilsk/retry/strategy"
)

sem := semaphore.New(5)

acquire := func(uint) error {
	return sem.Acquire(50 * time.Millisecond)
}

if err := retry.Retry(acquire, strategy.Limit(5), backoff.Linear(time.Second)); err != nil {
	// try again later
	return
}
defer sem.Release()

// do some heavy work
```

### Monitoring decorator

```go
type monitoredSemaphore struct {
	semaphore.Semaphore

	mu sync.Mutex
	timestamps []time.Time
}

func (sem *monitoredSemaphore) Acquire(timeout time.Duration) error {
	if err := sem.Semaphore.Acquire(timeout); err != nil {
		// trigger error counter to monitoring
		return err
	}
	sem.mu.Lock()
	sem.timestamps = append(sem.timestamps, time.Now())
	sem.mu.Unlock()
	// send current sem.Occupied() value to monitoring
	return nil
}

func (sem *monitoredSemaphore) Release() error {
	sem.mu.Lock()
	defer sem.mu.Unlock()
	if len(sem.timestamps) > 0 {
		var timestamp time.Time
		timestamp, sem.timestamps = sem.timestamps[0], sem.timestamps[1:]
		// send time.Since(timestamp) value to monitoring
	}
	return sem.Semaphore.Release()
}

func New(size int) semaphore.Semaphore {
	return &monitoredSemaphore{Semaphore: semaphore.New(size), timestamps: make([]time.Time, 0, size)}
}

sem := New(5)

if err := sem.Acquire(50*time.Millisecond); err != nil {
	// log error
	return
}
defer sem.Release()

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
