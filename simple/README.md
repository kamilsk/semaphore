> # semaphore/simple
>
> Simple non-blocking semaphore implementation with timeout based on channel.

[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore/simple?status.svg)](https://godoc.org/github.com/kamilsk/semaphore/simple)

## Usage

```go
sem := simple.New(5)

if err := sem.Acquire(50 * time.Millisecond); err != nil {
    // try again later
    return
}
defer sem.Release()

// do some heavy work
```

## Tips and tricks

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
