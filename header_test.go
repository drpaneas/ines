package ines_test

import (
	"testing"

	"github.com/drpaneas/ines"
)

func Test_hasHeader(t *testing.T) {
	t.Parallel()

	var (
		iNES2Header   = []byte{78, 69, 83, 26, 1, 1, 0, 8, 0, 0, 0, 0, 1, 0, 0, 1}
		iNES1Header   = []byte{78, 69, 83, 26, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		invalidHeader = []byte{79, 19, 23, 26, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	)

	type args struct {
		b []byte
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid ines 1.0 header",
			args: args{
				b: iNES1Header,
			},
			want: true,
		},
		{
			name: "valid ines 2.0 header",
			args: args{
				b: iNES2Header,
			},
			want: true,
		},
		{
			name: "invalid header",
			args: args{
				b: invalidHeader,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt2 := tt
		t.Run(tt2.name, func(t *testing.T) {
			t.Parallel()
			if got := ines.HasHeader(tt2.args.b); got != tt2.want {
				t.Errorf("HasHeader() = %v, want %v", got, tt2.want)
			}
		})
	}
}

func Test_isINES2(t *testing.T) {
	t.Parallel()

	var (
		iNES2Header = []byte{78, 69, 83, 26, 1, 1, 0, 8, 0, 0, 0, 0, 1, 0, 0, 1}
		iNES1Header = []byte{78, 69, 83, 26, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	)

	type args struct {
		b []byte
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid ines 1.0 header",
			args: args{
				b: iNES1Header,
			},
			want: false,
		},
		{
			name: "valid ines 2.0 header",
			args: args{
				b: iNES2Header,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt2 := tt
		t.Run(tt2.name, func(t *testing.T) {
			t.Parallel()
			if got := ines.IsINES2(tt2.args.b); got != tt2.want {
				t.Errorf("IsINES2() = %v, want %v", got, tt2.want)
			}
		})
	}
}
