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
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.1f at %d", tt.args.deg, tt.args.precision), func(t *testing.T) {
			if got := DirectionStr(tt.args.deg, tt.args.precision); got != tt.want {
				t.Errorf("DirectionStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
