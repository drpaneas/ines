package ines_test

import (
	"testing"

	"github.com/drpaneas/ines"
)

const wrongValue = 69

func Test_getTvSystemAndCPUPpuTiming(t *testing.T) {
	t.Parallel()

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
		tt2 := tt
		t.Run(tt2.name, func(t *testing.T) {
			t.Parallel()

			got, got1 := ines.GetTvSystemAndCPUPpuTiming(tt2.args.cpuPPUTiming)
			if got != tt2.wantTV {
				t.Errorf("getTvSystemAndCPUPpuTiming() got = %v, want %v", got, tt2.wantTV)
			}
			if got1 != tt2.wantCPU {
				t.Errorf("getTvSystemAndCPUPpuTiming() got1 = %v, want %v", got1, tt2.wantCPU)
			}
		})
	}
}
