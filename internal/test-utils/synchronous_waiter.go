package testutils

import (
	"sync"

	"github.com/thewizardplusplus/go-crawler/waiter"
)

// SynchronousWaiter ...
type SynchronousWaiter struct {
	actual sync.WaitGroup
	mock   waiter.Waiter
}

// NewSynchronousWaiter ...
func NewSynchronousWaiter(mock waiter.Waiter) *SynchronousWaiter {
	return &SynchronousWaiter{mock: mock}
}

// Add ...
func (waiter *SynchronousWaiter) Add(delta int) {
	waiter.actual.Add(delta)
	waiter.mock.Add(delta)
}

// Done ...
func (waiter *SynchronousWaiter) Done() {
	waiter.actual.Done()
	waiter.mock.Done()
}

// Wait ...
func (waiter *SynchronousWaiter) Wait() {
	waiter.actual.Wait()
}
