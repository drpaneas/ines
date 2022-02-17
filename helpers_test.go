package ines

import "testing"

func Test_readHighNibbleByte(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "01001111",
			args: args{b: 0b01001111},
			want: 0b0100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readHighNibbleByte(tt.args.b); got != tt.want {
				t.Errorf("readHighNibbleByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readLowNibbleByte(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "01000011",
			args: args{b: 0b01000011},
			want: 0b0011,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readLowNibbleByte(tt.args.b); got != tt.want {
				t.Errorf("readLowNibbleByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeNibbles(t *testing.T) {
	type args struct {
		highNibble byte
		lowNibble  byte
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "00101000",
			args: args{
				highNibble: readHighNibbleByte(0b00101000),
				lowNibble:  readLowNibbleByte(0b00101000),
			},
			want: 0b00101000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeNibbles(tt.args.highNibble, tt.args.lowNibble); got != tt.want {
				t.Errorf("mergeNibbles() = %v, want %v", got, tt.want)
			}
		})
	}
}
