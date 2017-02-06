package classic // import "github.com/kamilsk/semaphore/classic"

import (
	"sync"
)

// LockingSemaphore provides the functionality to limit bandwidth.
// https://en.wikipedia.org/wiki/Semaphore_(programming)#Operation_names
type LockingSemaphore interface {
	// P reduces the number of available slots for n.
	// It must be safe to call P concurrently on a single semaphore.
	P(n int)
	// V increases the number of available slots for n.
	// It must be safe to call V concurrently on a single semaphore.
	V(n int)
}

// NewLocking constructs a new thread-safe LockingSemaphore with the given capacity.
func NewLocking(capacity int) LockingSemaphore {
	return make(semaphore, capacity)
}

func (sem semaphore) P(n int) {
	for i := 0; i < n; i++ {
		sem <- struct{}{}
	}
}

func (sem semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-sem
	}
}

// SyncingSemaphore provides the functionality to synchronize multiple goroutines.
// https://en.wikipedia.org/wiki/Semaphore_(programming)#Semantics_and_implementation
type SyncingSemaphore interface {
	// Signal reports about completion of goroutine work.
	// The method increases the number of available slots for 1.
	// It must be safe to call Signal concurrently on a single semaphore.
	Signal()
	// Wait starts to wait n goroutines.
	// The method reduces the number of available slots for n.
	// It must be safe to call Wait concurrently on a single semaphore.
	Wait(n int)
}

// NewSyncing constructs a new full-filled thread-safe SyncingSemaphore with the given capacity.
func NewSyncing(capacity int) SyncingSemaphore {
	sem := make(semaphore, capacity)
	sem.P(capacity)
	return sem
}

func (sem semaphore) Signal() {
	sem.V(1)
}

func (sem semaphore) Wait(n int) {
	sem.P(n)
}

// BinarySemaphore represents the classic binary semaphore with mutex-like interface.
type BinarySemaphore interface {
	sync.Locker
}

// NewBinary constructs a new thread-safe BinarySemaphore with capacity equals to one.
func NewBinary() BinarySemaphore {
	return make(semaphore, 1)
}

func (sem semaphore) Lock() {
	sem.P(1)
}

func (sem semaphore) Unlock() {
	sem.V(1)
}

type semaphore chan struct{}
