package libwx

import (
	"fmt"
	"testing"
)

func TestDirectionStr(t *testing.T) {
	type args struct {
		deg       Degree
		precision DirectionStrPrecision
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{deg: 350, precision: DirectionStrPrecision1},
			want: "N",
		},
		{
			args: args{deg: 230, precision: DirectionStrPrecision2},
			want: "SW",
		},
		{
			args: args{deg: 165, precision: DirectionStrPrecision3},
			want: "SSE",
		},
		{
			args: args{deg: 0, precision: DirectionStrPrecision3},
			want: "N",
		},
		{
			args: args{deg: 360, precision: DirectionStrPrecision3},
			want: "N",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.1f at %d", tt.args.deg, tt.args.precision), func(t *testing.T) {
			if got := DirectionStr(tt.args.deg, tt.args.precision); got != tt.want {
				t.Errorf("DirectionStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClampedDegree(t *testing.T) {
	tests := []struct {
		in   float64
		want Degree
	}{
		{in: 0, want: 360},
		{in: 360, want: 360},
		{in: 361, want: 1},
		{in: 0.1, want: 0.1},
		{in: -1, want: 359},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.1f", tt.in), func(t *testing.T) {
			if got := ClampedDegree(tt.in); got != tt.want {
				t.Errorf("ClampedDegree() = %v, want %v", got, tt.want)
			}
		})
	}
}
