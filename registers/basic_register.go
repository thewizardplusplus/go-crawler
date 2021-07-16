package registers

import (
	"sync"
)

// BasicRegister ...
type BasicRegister struct {
	registeredValues *sync.Map
}

// NewBasicRegister ...
func NewBasicRegister() BasicRegister {
	return BasicRegister{
		registeredValues: new(sync.Map),
	}
}
