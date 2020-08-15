package waiter

// Waiter ...
type Waiter interface {
	Add(delta int)
	Done()
}
