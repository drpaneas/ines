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
		tt2 := tt
		t.Run(tt2.name, func(t *testing.T) {
			if got := readHighNibbleByte(tt2.args.b); got != tt2.want {
				t.Errorf("readHighNibbleByte() = %v, want %v", got, tt2.want)
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
		tt2 := tt
		t.Run(tt2.name, func(t *testing.T) {
			if got := readLowNibbleByte(tt2.args.b); got != tt2.want {
				t.Errorf("readLowNibbleByte() = %v, want %v", got, tt2.want)
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
		tt2 := tt
		t.Run(tt2.name, func(t *testing.T) {
			if got := mergeNibbles(tt2.args.highNibble, tt2.args.lowNibble); got != tt2.want {
				t.Errorf("mergeNibbles() = %v, want %v", got, tt2.want)
			}
		})
	}
}
