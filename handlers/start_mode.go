package handlers

import (
	"sync"
)

type startMode int

const (
	notStarted startMode = iota
	started
	startedConcurrently
)

type startModeHolder struct {
	lock sync.RWMutex
	mode startMode
}

func (holder *startModeHolder) getStartMode() startMode {
	holder.lock.RLock()
	defer holder.lock.RUnlock()

	return holder.mode
}

func (holder *startModeHolder) setStartModeOnce(mode startMode) {
	holder.lock.Lock()
	defer holder.lock.Unlock()

	if holder.mode == notStarted {
		holder.mode = mode
	}
}
