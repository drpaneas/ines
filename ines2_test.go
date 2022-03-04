package ines

import "testing"

const wrongValue = 69

func Test_getTvSystemAndCPUPpuTiming(t *testing.T) {
	type args struct {
		cpuPPUTiming int
	}
	tests := []struct {
		name    string
		args    args
		wantTV  string
		wantCPU string
	}{
		{
			name:    "default values",
			args:    args{cpuPPUTiming: wrongValue},
			wantTV:  "Unknown",
			wantCPU: "Unknown",
		},
		{
			name:    "NASA Region",
			args:    args{cpuPPUTiming: 0},
			wantTV:  "North America, Japan, South Korea, Taiwan",
			wantCPU: "RP2C02 (\"NTSC NES\")",
		},
		{
			name:    "EMEA Region",
			args:    args{cpuPPUTiming: 1},
			wantTV:  "Western Europe, Australia",
			wantCPU: "RP2C07 (\"Licensed PAL NES\")",
		},
		{
			name:    "APAC Region",
			args:    args{cpuPPUTiming: 3},
			wantTV:  "Eastern Europe, Russia, Mainland China, India, Africa",
			wantCPU: "UMC 6527P (\"Dendy\")",
		},
		{
			name:    "Same For All Regions",
			args:    args{cpuPPUTiming: 2},
			wantTV:  "Identical ROM content in both NTSC and PAL countries",
			wantCPU: "Multiple-region",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getTvSystemAndCPUPpuTiming(tt.args.cpuPPUTiming)
			if got != tt.wantTV {
				t.Errorf("getTvSystemAndCPUPpuTiming() got = %v, want %v", got, tt.wantTV)
			}
			if got1 != tt.wantCPU {
				t.Errorf("getTvSystemAndCPUPpuTiming() got1 = %v, want %v", got1, tt.wantCPU)
			}
		})
	}
}
