package hl7v2

import (
	"reflect"
	"testing"
)

func TestParseDelimiters(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Delimiters
		wantErr bool
	}{
		{
			name: "test_default",
			args: args{
				b: []byte("|^~\\&"),
			},
			want:    DefaultDelimiters(),
			wantErr: false,
		},
		{
			name: "test_invalid_length",
			args: args{
				b: []byte("|^~"),
			},
			want:    DefaultDelimiters(),
			wantErr: true,
		},
		{
			name: "test_invalid_char",
			args: args{
				b: []byte("|^~\\|"),
			},
			want:    DefaultDelimiters(),
			wantErr: true,
		},
		{
			name: "test_truncation",
			args: args{
				b: []byte("|^~\\&!"),
			},
			want:    NewDelimiters(SetTruncation('!')),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDelimiters(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDelimiters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDelimiters() = %v, want %v", got, tt.want)
			}
		})
	}
}
