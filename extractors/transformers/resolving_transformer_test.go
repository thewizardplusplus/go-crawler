package transformers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_selectBaseTag(test *testing.T) {
	type args struct {
		data []byte
	}

	for _, data := range []struct {
		name string
		args args
		want string
	}{
		{
			name: "without the base tag",
			args: args{
				data: []byte(`
					<ul>
						<li><a href="http://example.com/1">1</a></li>
						<li><a href="http://example.com/2">2</a></li>
					</ul>
				`),
			},
			want: "",
		},
		{
			name: "with the base tag without the href attribute",
			args: args{
				data: []byte(`
					<base target="_blank" />

					<ul>
						<li><a href="http://example.com/1">1</a></li>
						<li><a href="http://example.com/2">2</a></li>
					</ul>
				`),
			},
			want: "",
		},
		{
			name: "with the base tag with the href attribute",
			args: args{
				data: []byte(`
					<base href="http://example.com/" />

					<ul>
						<li><a href="1">1</a></li>
						<li><a href="2">2</a></li>
					</ul>
				`),
			},
			want: "http://example.com/",
		},
		{
			name: "with the several base tags with the href attribute",
			args: args{
				data: []byte(`
					<base href="http://example.com/1/" />
					<base href="http://example.com/2/" />

					<ul>
						<li><a href="3">3</a></li>
						<li><a href="4">4</a></li>
					</ul>
				`),
			},
			want: "http://example.com/1/",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := selectBaseTag(data.args.data)

			assert.Equal(test, data.want, got)
		})
	}
}
