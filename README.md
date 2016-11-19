> # Semaphore
>
> Simple non-blocking semaphore implementation with timeout written on Go.

[![Build Status](https://travis-ci.org/kamilsk/semaphore.svg?branch=master)](https://travis-ci.org/kamilsk/semaphore)
[![Coverage Status](https://coveralls.io/repos/github/kamilsk/semaphore/badge.svg)](https://coveralls.io/github/kamilsk/semaphore)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamilsk/semaphore)](https://goreportcard.com/report/github.com/kamilsk/semaphore)
[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore?status.svg)](https://godoc.org/github.com/kamilsk/semaphore)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](LICENSE.md)

## Usage

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

### HTTP request limiter

```go
limiter := semaphore.New(10)

http.Handle("/do-some-heavy-work", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	if err := limiter.Acquire(time.Second); err != nil {
		// try again after 1 minute
		rw.Header().Add("Retry-After", time.Minute.String())
		rw.WriteHeader(http.StatusServiceUnavailable)
		rw.Write(nil)
		return
	}

	// do some heavy work
}))
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

### Retry to acquire a few times with line breaks between attempts

```go
import (
	"github.com/kamilsk/retrier"
	"github.com/kamilsk/retrier/backoff"
	"github.com/kamilsk/retrier/strategy"
)

sem := semaphore.New(5)

acquire := func(uint) error {
	return sem.Acquire(50 * time.Millisecond)
}

if err := retrier.Retry(acquire, strategy.Limit(5), backoff.Linear(time.Second)); err != nil {
	// try again later
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

## Notes

### Articles

- [Semaphore Pattern](http://tmrts.com/go-patterns/synchronization/semaphore.html)
- [Semaphores - Go Language Patterns](https://sites.google.com/site/gopatterns/concurrency/semaphores)
- [Semaphore (programming)](https://en.wikipedia.org/wiki/Semaphore_(programming))

### Examples

- [github.com/goinaction/code/chapter7/patterns/semaphore](https://github.com/goinaction/code/tree/master/chapter7/patterns/semaphore)

### Another implementations

#### Similar

- [github.com/eapache/go-resiliency/semaphore.Semaphore](https://github.com/eapache/go-resiliency/blob/master/semaphore/semaphore.go)
- [github.com/dropbox/godropbox/sync2.Semaphore](https://github.com/dropbox/godropbox/blob/master/sync2/semaphore.go)
- [github.com/tarndt/sema.TimeoutCountingSema](https://github.com/tarndt/sema/blob/master/sema.go)
- [github.com/abiosoft/semaphore.Semaphore](https://github.com/abiosoft/semaphore/blob/master/semaphore.go)
- [github.com/vada-ir/semaphore.Semaphore](https://github.com/vada-ir/semaphore/blob/master/semaphore.go)
- [github.com/avezila/psem.Sem](https://github.com/avezila/psem/blob/master/psem.go)
- [github.com/jsipprell/go-semaphore.WaitableSemaphore](https://github.com/jsipprell/go-semaphore/blob/master/semaphore.go)

#### Others

##### Not locking

- [github.com/toolkits/concurrent/semaphore](https://github.com/toolkits/concurrent/tree/master/semaphore)
- [github.com/dexterous/semaphore](https://github.com/dexterous/semaphore)
- [github.com/seanjohnno/semaphore](https://github.com/seanjohnno/semaphore)
- [github.com/tmthrgd/go-sem](https://github.com/tmthrgd/go-sem)

##### Locking

- [github.com/tj/go-sync/semaphore](https://github.com/tj/go-sync/tree/master/semaphore)
- [github.com/carlmjohnson/go-utils/semaphore](https://github.com/carlmjohnson/go-utils/tree/master/semaphore)
- [github.com/pivotal-golang/semaphore](https://github.com/pivotal-golang/semaphore)
- [github.com/andreyvit/sem](https://github.com/andreyvit/sem)
- [github.com/spektroskop/semaphore](https://github.com/spektroskop/semaphore)
- [github.com/opencoff/go-lib/sem](https://github.com/opencoff/go-lib/tree/master/sem)
- [github.com/nicholasjackson/bench/semaphore](https://github.com/nicholasjackson/bench/tree/master/semaphore)
- [github.com/cognusion/semaphore](https://github.com/cognusion/semaphore)
- [github.com/riobard/go-semaphore](https://github.com/riobard/go-semaphore)
