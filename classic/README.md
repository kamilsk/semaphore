> # semaphore/classic
>
> Classic Semaphore pattern implementation based on channel.

[![GoDoc](https://godoc.org/github.com/kamilsk/semaphore/classic?status.svg)](https://godoc.org/github.com/kamilsk/semaphore/classic)

## Usage

### Locking

```go
sem := classic.NewLocking(2)

for i := 0; i < 2; i++ {
    go func() {
        defer sem.P(1)
        fmt.Println("work is done")
    }()
}

sem.V(2)
fmt.Println("all work is done")
```

### Syncing

```go
sem := classic.NewSyncing(2)

for i := 0; i < 2; i++ {
    go func() {
        defer sem.Signal()
        fmt.Println("process is finished")
    }()
}

sem.Wait(2)
fmt.Println("all processes are finished")
```

### Binary

```go
binary := classic.NewBinary()

var shared string

go func() {
    binary.Lock()
    shared = "a"
    binary.Unlock()
}()

// just enough to yield the scheduler and let the goroutines work off
time.Sleep(time.Millisecond)

binary.Lock()
shared = "b"
binary.Unlock()
```
