package semaphore

import "runtime"

var def = New(runtime.GOMAXPROCS(0))

// Capacity returns a capacity of the default semaphore.
func Capacity() int {
	return def.Capacity()
}

// Occupied returns a current number of occupied slots of the default semaphore.
func Occupied() int {
	return def.Occupied()
}

// Release releases the previously occupied slot of the default semaphore.
func Release() error {
	return def.Release()
}

// Acquire tries to reduce the number of available slots of the default semaphore for 1.
// The operation can be canceled using deadline channel. In this case,
// an appropriate error will be returned.
func Acquire(deadline <-chan struct{}) (ReleaseFunc, error) {
	return def.Acquire(deadline)
}
