package registers

import (
	"context"
	"sync"
	"testing"
	"testing/iotest"

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
		{
			name: "success with an unregistered value",
			fields: fields{
				registeredValues: new(sync.Map),
			},
			args: args{
				ctx: context.Background(),
				key: "key",
				registeringHandler: func() RegisteringHandler {
					registeringHandler := new(MockRegisteringHandler)
					registeringHandler.
						On("HandleRegistering", context.Background(), "key").
						Return("value", nil)

					return registeringHandler
				}(),
			},
			wantValue: "value",
			wantErr:   assert.NoError,
		},
		{
			name: "success with a registered value",
			fields: fields{
				registeredValues: func() *sync.Map {
					registeredRobotsTXT := new(sync.Map)
					registeredRobotsTXT.Store("key", "value")

					return registeredRobotsTXT
				}(),
			},
			args: args{
				ctx:                context.Background(),
				key:                "key",
				registeringHandler: new(MockRegisteringHandler),
			},
			wantValue: "value",
			wantErr:   assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				registeredValues: new(sync.Map),
			},
			args: args{
				ctx: context.Background(),
				key: "key",
				registeringHandler: func() RegisteringHandler {
					registeringHandler := new(MockRegisteringHandler)
					registeringHandler.
						On("HandleRegistering", context.Background(), "key").
						Return(nil, iotest.ErrTimeout)

					return registeringHandler
				}(),
			},
			wantValue: nil,
			wantErr:   assert.Error,
		},
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
