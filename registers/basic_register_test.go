package registers

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBasicRegister(test *testing.T) {
	got := NewBasicRegister()

	assert.Equal(test, new(sync.Map), got.registeredValues)
}

func TestBasicRegister_RegisterValue(test *testing.T) {
	type fields struct {
		registeredValues *sync.Map
	}
	type args struct {
		ctx                context.Context
		key                interface{}
		registeringHandler RegisteringHandler
	}

	for _, data := range []struct {
		name      string
		fields    fields
		args      args
		wantValue interface{}
		wantErr   assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(t *testing.T) {
			register := BasicRegister{
				registeredValues: data.fields.registeredValues,
			}
			gotValue, gotErr := register.RegisterValue(
				data.args.ctx,
				data.args.key,
				data.args.registeringHandler.HandleRegistering,
			)

			mock.AssertExpectationsForObjects(test, data.args.registeringHandler)
			assert.Equal(t, data.wantValue, gotValue)
			data.wantErr(test, gotErr)
		})
	}
}
