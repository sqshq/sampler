package gauge

import "testing"

func Test_calculatePercent(t *testing.T) {
	type args struct {
		cur, min, max float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"should calculate percent between 0 and 60", args{cur: 45, min: 0, max: 60}, 75},
		{"should calculate percent between 10 and 60", args{cur: 59, min: 10, max: 60}, 98},
		{"should calculate percent between -20 and 60", args{cur: 0, min: -20, max: 60}, 25},
		{"should calculate percent when cur value = min value", args{cur: -10, min: -10, max: 60}, 0},
		{"should calculate percent when min value = max value", args{cur: -124, min: -124, max: -124}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePercent(tt.args.cur, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("calculatePercent() = %v, want %v", got, tt.want)
			}
		})
	}
}
