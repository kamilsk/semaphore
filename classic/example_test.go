package classic_test

import (
	"fmt"
	"time"

	"github.com/kamilsk/semaphore/classic"
)

// This example shows how to work with locking semaphore.
func Example_lockingSemaphore() {
	sem := classic.NewLocking(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer sem.P(1)
			fmt.Println("work is done")
		}()
	}

	sem.V(2)
	fmt.Println("all work is done")

	// Output:
	// work is done
	// work is done
	// all work is done
}

// This example shows how to work with syncing semaphore.
func Example_syncingSemaphore() {
	sem := classic.NewSyncing(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer sem.Signal()
			fmt.Println("process is finished")
		}()
	}

	sem.Wait(2)
	fmt.Println("all processes are finished")

	// Output:
	// process is finished
	// process is finished
	// all processes are finished
}

// This example shows hot to work with binary semaphore.
func Example_binarySemaphore() {
	binary := classic.NewBinary()

	var shared string

	go func() {
		binary.Lock()
		defer binary.Unlock()
		shared = "a"
	}()

	// just enough to yield the scheduler and let the goroutines work off
	time.Sleep(time.Millisecond)

	binary.Lock()
	defer binary.Unlock()
	shared = "b"

	fmt.Printf("shared value is equals to %q", shared)

	// Output: shared value is equals to "b"
}
