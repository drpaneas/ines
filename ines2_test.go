package ines

import "testing"

func Test_getTvSystemAndCpuPpuTiming(t *testing.T) {
	type args struct {
		cpuPPUTiming int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getTvSystemAndCpuPpuTiming(tt.args.cpuPPUTiming)
			if got != tt.want {
				t.Errorf("getTvSystemAndCpuPpuTiming() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getTvSystemAndCpuPpuTiming() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
