package semaphore

import "sync/atomic"

type draft struct {
	state    uint32
	capacity uint32
}

func (semaphore *draft) Release() error {
	panic("implement me")
}

func (semaphore *draft) Acquire(breaker BreakCloser, places ...uint32) (Releaser, error) {
	_ = reduce(places...)
	panic("implement me")
}

func (semaphore *draft) Try(breaker Breaker, places ...uint32) (Releaser, error) {
	_ = reduce(places...)
	panic("implement me")
}

func (semaphore *draft) Signal(breaker Breaker) <-chan Releaser {
	panic("implement me")
}

func (semaphore *draft) Peek() uint32 {
	return atomic.LoadUint32(&semaphore.state)
}

func (semaphore *draft) Size(new uint32) uint32 {
	current := atomic.LoadUint32(&semaphore.capacity)
	if new != 0 {
		atomic.StoreUint32(&semaphore.capacity, new)
	}
	return current
}

func reduce(places ...uint32) uint32 {
	var capacity uint32
	for _, size := range places {
		capacity += size
	}
	if capacity == 0 {
		return 1
	}
	return capacity
}
