package examples_test

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kamilsk/semaphore"
)

type Pool struct {
	sem  semaphore.Semaphore
	work chan func()
}

func (p *Pool) Schedule(task func()) {
	select {
	case p.work <- task:
	case release := <-p.sem.Signal(nil):
		go p.worker(task, release)
	}
}

func (p *Pool) worker(task func(), release semaphore.ReleaseFunc) {
	defer release()
	var ok bool
	for {
		task()
		task, ok = <-p.work
		if !ok {
			return
		}
	}
}

func New(size int) *Pool {
	return &Pool{
		sem:  semaphore.New(size),
		work: make(chan func()),
	}
}

// This example shows how to create a pool of workers based on the semaphore.
func Example_poolOfWorkers() {
	var ok, fail int32 = 0, 5

	wg := &sync.WaitGroup{}
	do := func() {
		atomic.AddInt32(&ok, 1)
		atomic.AddInt32(&fail, -1)
		wg.Done()
	}
	pool := New(int(fail / 2))

	wg.Add(int(fail))
	for i, total := 0, int(fail); i < total; i++ {
		pool.Schedule(do)
	}
	wg.Wait()

	fmt.Printf("success: %d, failure: %d \n", ok, fail)
	// Output: success: 5, failure: 0
}
