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
			fmt.Println("work done")
			sem.P(1)
		}()
	}

	sem.V(2)
	fmt.Println("all works done")

	// Output:
	// work done
	// work done
	// all works done
}

// This example shows how to work with syncing semaphore.
func Example_syncingSemaphore() {
	sem := classic.NewSyncing(2)

	for i := 0; i < 2; i++ {
		go func() {
			fmt.Println("process finished")
			sem.Signal()
		}()
	}

	sem.Wait(2)
	fmt.Println("all processes finished")

	// Output:
	// process finished
	// process finished
	// all processes finished
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
