package registers

import (
	"sync"
)

// BasicRegister ...
type BasicRegister struct {
	registeredValues *sync.Map
}
