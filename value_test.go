package hl7v2

import (
	"reflect"
	"testing"
)

func TestUnescapeValue(t *testing.T) {
	type args struct {
		v      []byte
		delims *Delimiters
	}
	tests := []struct {
		name string
		args args
		want Value
	}{
		{
			name: "test sub",
			args: args{
				v:      []byte("testing \\T\\ stuff"),
				delims: DefaultDelimiters(),
			},
			want: NewStringValue("testing & stuff"),
		},
		{
			name: "test field",
			args: args{
				v:      []byte("testing \\F\\ stuff"),
				delims: DefaultDelimiters(),
			},
			want: NewStringValue("testing | stuff"),
		},
		{
			name: "test rep",
			args: args{
				v:      []byte("testing \\R\\ stuff"),
				delims: DefaultDelimiters(),
			},
			want: NewStringValue("testing ~ stuff"),
		},
		{
			name: "test escape",
			args: args{
				v:      []byte("testing \\E\\ stuff"),
				delims: DefaultDelimiters(),
			},
			want: NewStringValue("testing \\ stuff"),
		},
		{
			name: "test comp",
			args: args{
				v:      []byte("testing \\S\\ stuff"),
				delims: DefaultDelimiters(),
			},
			want: NewStringValue("testing ^ stuff"),
		},
		{
			name: "test lf",
			args: args{
				v:      []byte("testing \\X0A\\ stuff"),
				delims: DefaultDelimiters(),
			},
			want: NewStringValue("testing \n stuff"),
		},
		{
			name: "test cr",
			args: args{
				v:      []byte("testing \\X0D\\ stuff"),
				delims: DefaultDelimiters(),
			},
			want: NewStringValue("testing \r stuff"),
		},
		{
			name: "test br",
			args: args{
				v:      []byte("testing \\.br\\ stuff"),
				delims: DefaultDelimiters(),
			},
			want: NewStringValue("testing \r stuff"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnescapeValue(tt.args.v, tt.args.delims); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnescapeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
