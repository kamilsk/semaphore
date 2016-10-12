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
	timers []time.Time
}

func (sem *monitoredSemaphore) Acquire(timeout time.Duration) error {
	if err := sem.Semaphore.Acquire(timeout); err != nil {
		// trigger error counter to monitoring
		return err
	}
	// send current sem.Occupied() value to monitoring
	sem.timers = append(sem.timers, time.Now())
	return nil
}

func (sem *monitoredSemaphore) Release() error {
	sem.mu.Lock()
	defer sem.mu.Unlock()
	if len(sem.timers) > 0 {
		timer := sem.timers[0]
		sem.timers := sem.timers[1:]
		// send time.Since(timer) value to monitoring
	}
	return sem.Semaphore.Release()
}

func New(size int) semaphore.Semaphore {
	return &monitoredSemaphore{Semaphore: semaphore.New(5)}
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
