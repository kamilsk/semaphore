package classic

import (
	"sync"
)

// BlockingSemaphore ...
// https://en.wikipedia.org/wiki/Semaphore_(programming)#Operation_names
type BlockingSemaphore interface {
	// P ...
	P(n int)
	// V ...
	V(n int)
}

// NewBlockingSemaphore constructs a new BlockingSemaphore with the given number of places.
func NewBlockingSemaphore(size int) BlockingSemaphore {
	return make(semaphore, size)
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

// ProcessSemaphore ...
// https://en.wikipedia.org/wiki/Semaphore_(programming)#Semantics_and_implementation
type ProcessSemaphore interface {
	Signal()
	Wait(n int)
}

// NewProcessSemaphore constructs a new ProcessSemaphore with the given number of places.
func NewProcessSemaphore(size int) ProcessSemaphore {
	sem := make(semaphore, size)
	sem.P(size)
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

// NewBinarySemaphore constructs a new BinarySemaphore with one place.
func NewBinarySemaphore() BinarySemaphore {
	return make(semaphore, 1)
}

func (sem semaphore) Lock() {
	sem.P(1)
}

func (sem semaphore) Unlock() {
	sem.V(1)
}

type semaphore chan struct{}
