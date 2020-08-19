package testutils

import (
	"sync"

	"github.com/thewizardplusplus/go-crawler/waiter"
)

// SynchronousWaiter ...
type SynchronousWaiter struct {
	mock   waiter.Waiter
	actual sync.WaitGroup
}

// NewSynchronousWaiter ...
func NewSynchronousWaiter(mock waiter.Waiter) *SynchronousWaiter {
	return &SynchronousWaiter{mock: mock}
}

// Add ...
func (waiter *SynchronousWaiter) Add(delta int) {
	waiter.mock.Add(delta)
	waiter.actual.Add(delta)
}

// Done ...
func (waiter *SynchronousWaiter) Done() {
	waiter.mock.Done()
	waiter.actual.Done()
}

// Wait ...
func (waiter *SynchronousWaiter) Wait() {
	waiter.actual.Wait()
}
