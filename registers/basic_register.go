package registers

import (
	"context"
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

// RegisterValue ...
func (register BasicRegister) RegisterValue(
	ctx context.Context,
	key interface{},
	registeringHandler func(ctx context.Context, key interface{}) (
		value interface{},
		err error,
	),
) (
	value interface{},
	err error,
) {
	value, ok := register.registeredValues.Load(key)
	if !ok {
		var err error
		value, err = registeringHandler(ctx, key)
		if err != nil {
			return nil, err
		}

		register.registeredValues.Store(key, value)
	}

	return value, nil
}
