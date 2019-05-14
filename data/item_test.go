package data

import "testing"

func Test_cleanupOutput(t *testing.T) {
	type args struct {
		output string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should trim everything before carriage return", args{">>\rtext"}, "text"},
		{"should trim carriage return at the end", args{"text\r"}, "text"},
		{"should remove tabs and spaces", args{"\t\t\ntext "}, "text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanupOutput(tt.args.output); got != tt.want {
				t.Errorf("cleanupOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
