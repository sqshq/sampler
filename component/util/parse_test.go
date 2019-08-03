package util

import "testing"

func TestParseFloat(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{"should parse a regular number", args{"123"}, 123, false},
		{"should parse a float with decimal point", args{"123.456"}, 123.456, false},
		{"should parse a float with decimal comma", args{"123,456"}, 123.456, false},
		{"should parse a regular number with spaces, tabs and line breaks", args{"         \t 123 \t \n    "}, 123, false},
		{"should parse a last line in the given string", args{"123\n456"}, 456, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFloat(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
