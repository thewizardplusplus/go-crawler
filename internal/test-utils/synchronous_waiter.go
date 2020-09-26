package testutils

import (
	"sync"

	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

// SynchronousWaiter ...
type SynchronousWaiter struct {
	mock   syncutils.WaitGroup
	actual sync.WaitGroup
}

// NewSynchronousWaiter ...
func NewSynchronousWaiter(mock syncutils.WaitGroup) *SynchronousWaiter {
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
