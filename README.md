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

## Tips and tricks

### Monitoring decorator

```go
type MonitoredSemaphore struct {
	Semaphore

	mu sync.Mutex
	timers []time.Time
}

func (sem *MonitoredSemaphore) Acquire(timeout time.Duration) error {
	if err := sem.Semaphore.Acquire(timeout); err != nil {
		// trigger error counter in monitoring
		return err
	}
	// send current sem.Occupied() value to monitoring
	sem.timers = append(sem.timers, time.Now())
	return nil
}

func (sem *MonitoredSemaphore) Release() {
	sem.mu.Lock()
	defer sem.mu.Unlock()
	if len(sem.timers) > 0 {
		timer := sem.timers[0]
		sem.timers := sem.timers[1:]
		// send time.Since(timer) value to monitoring
	}
	sem.Semaphore.Release()
}

var sem Semaphore = &MonitoredSemaphore{Semaphore: new semaphore.New(5)}

if err := sem.Acquire(50*time.Millisecond); err != nil {
	// log error
}
defer sem.Release()
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
