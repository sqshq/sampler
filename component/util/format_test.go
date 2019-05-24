package util

import "testing"

func TestFormatValue(t *testing.T) {
	type args struct {
		value float64
		scale int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should format float value with scale = 1", args{94.123, 1}, "94.1"},
		{"should format float value with scale = 0", args{94.123, 0}, "94"},
		{"should format float value with trailing zeros", args{94.100, 5}, "94.1"},
		{"should format float value with radix char and trailing zeros", args{9423000.00123, 2}, "9,423,000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatValue(tt.args.value, tt.args.scale); got != tt.want {
				t.Errorf("FormatValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
