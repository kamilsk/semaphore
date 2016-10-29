package classic_test

import (
	"fmt"
	"time"

	"github.com/kamilsk/semaphore/classic"
)

// This example shows how to work with blocking semaphore.
func Example_blockingSemaphore() {
	semaphore := classic.NewBlockingSemaphore(2)

	for i := 0; i < 2; i++ {
		go func() {
			fmt.Println("work is done")
			semaphore.P(1)
		}()
	}

	semaphore.V(2)
	fmt.Println("all work is done")

	// Output:
	// work is done
	// work is done
	// all work is done
}

// This example shows how to work with process semaphore.
func Example_processSemaphore() {
	process := classic.NewProcessSemaphore(2)

	for i := 0; i < 2; i++ {
		go func() {
			fmt.Println("process has finished")
			process.Signal()
		}()
	}

	process.Wait(2)
	fmt.Println("all processes have finished")

	// Output:
	// process has finished
	// process has finished
	// all processes have finished
}

// This example shows hot to work with binary semaphore.
func Example_binarySemaphore() {
	binary := classic.NewBinarySemaphore()
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

	fmt.Printf("shared value is %q", shared)
	// Output: shared value is "b"
}
