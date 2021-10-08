package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_startModeHolder_getStartMode(test *testing.T) {
	holder := &startModeHolder{mode: 23}
	receivedStartMode := holder.getStartMode()

	assert.Equal(test, startMode(23), receivedStartMode)
}

func Test_startModeHolder_setStartModeOnce(test *testing.T) {
	type fields struct {
		mode startMode
	}
	type args struct {
		mode startMode
	}

	for _, data := range []struct {
		name            string
		fields          fields
		args            args
		wantedStartMode startMode
	}{
		{
			name: "with the default start mode",
			fields: fields{
				mode: notStarted,
			},
			args: args{
				mode: 42,
			},
			wantedStartMode: 42,
		},
		{
			name: "with the not default start mode",
			fields: fields{
				mode: 23,
			},
			args: args{
				mode: 42,
			},
			wantedStartMode: 23,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			holder := &startModeHolder{
				mode: data.fields.mode,
			}
			holder.setStartModeOnce(data.args.mode)

			assert.Equal(test, data.wantedStartMode, holder.getStartMode())
		})
	}
}
