package gauge

import "testing"

func Test_calculatePercent(t *testing.T) {
	type args struct {
		g *Gauge
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"should calculate percent between 0 and 60", args{&Gauge{minValue: 0, maxValue: 60, curValue: 45}}, 75},
		{"should calculate percent between 10 and 60", args{&Gauge{minValue: 10, maxValue: 60, curValue: 59}}, 98},
		{"should calculate percent between -20 and 60", args{&Gauge{minValue: -20, maxValue: 60, curValue: 0}}, 25},
		{"should calculate percent when cur value = min value", args{&Gauge{minValue: -10, maxValue: 60, curValue: -10}}, 0},
		{"should calculate percent when min value = max value", args{&Gauge{minValue: -124, maxValue: -124, curValue: -124}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePercent(tt.args.g); got != tt.want {
				t.Errorf("calculatePercent() = %v, want %v", got, tt.want)
			}
		})
	}
}
