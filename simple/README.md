> # semaphore/simple
>
> Semaphore pattern implementation with timeout of lock/unlock operations based on channel.

[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore/simple?status.svg)](https://godoc.org/github.com/kamilsk/semaphore/simple)

## Usage

### Basic

```go
sem := simple.New(5)

if err := sem.Acquire(50 * time.Millisecond); err != nil {
    // try again later
    return
}
defer sem.Release()

// do some heavy work
```

### HTTP request limiter

```go
limiter := simple.New(10)

http.Handle("/do-some-heavy-work", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
    if err := limiter.Acquire(time.Second); err != nil {
        // try again after 1 minute
        rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
        rw.Header().Set("Retry-After", time.Minute.String())
        rw.Header().Set("X-Content-Type-Options", "nosniff")
        rw.WriteHeader(http.StatusServiceUnavailable)
        return
    }

    // do some heavy work
}))
```

## Tips and tricks

### Retry to acquire a few times with line breaks between attempts

```go
import (
    "github.com/kamilsk/retrier"
    "github.com/kamilsk/retrier/backoff"
    "github.com/kamilsk/retrier/strategy"
)

sem := simple.New(5)

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

### Monitoring decorator

```go
type monitoredSemaphore struct {
    simple.Semaphore

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
    if len(sem.timestamps) > 0 {
        var timestamp time.Time
        timestamp, sem.timestamps = sem.timestamps[0], sem.timestamps[1:]
        sem.mu.Unlock()
        // send time.Since(timestamp) value to monitoring
    }
    sem.mu.Unlock()
    return sem.Semaphore.Release()
}

func New(capacity int) simple.Semaphore {
    return &monitoredSemaphore{Semaphore: simple.New(capacity), timestamps: make([]time.Time, 0, capacity)}
}

sem := New(5)

if err := sem.Acquire(50 * time.Millisecond); err != nil {
    // log error and try again later
    return
}
defer sem.Release()

// do some heavy work
```
