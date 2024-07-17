package libwx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO(cdzombak): testing for calculations in libwx.go

func Test_WetBulbC(t *testing.T) {
	r := require.New(t)
	eq := CurriedFloat64Equal(Tolerance0)

	cases := []struct {
		t        TempC
		rh       RelHumidity
		expected TempC
	}{
		{TempC(-19), RelHumidity(80), TempC(-20)},
		{TempC(-14), RelHumidity(90), TempC(-15)},
		{TempC(-1), RelHumidity(30), TempC(-5)},
		{TempC(-1), RelHumidity(40), TempC(-5)},
		{TempC(-3.5), RelHumidity(90), TempC(-5)},
		{TempC(5), RelHumidity(30), TempC(0)},
		{TempC(1), RelHumidity(90), TempC(0)},
		{TempC(6.5), RelHumidity(90), TempC(5)},
		{TempC(10), RelHumidity(50), TempC(5)},
		{TempC(19), RelHumidity(5), TempC(5)},
		{TempC(10), RelHumidity(99), TempC(10)},
		{TempC(25), RelHumidity(10), TempC(10)},
		{TempC(15), RelHumidity(50), TempC(10)},
		{TempC(20), RelHumidity(60), TempC(15)},
		{TempC(20), RelHumidity(99), TempC(20)},
		{TempC(40), RelHumidity(70), TempC(35)},
		{TempC(40), RelHumidity(99), TempC(40)},
		{TempC(50), RelHumidity(11), TempC(25)},
		{TempC(44), RelHumidity(50), TempC(35)},
		{TempC(46), RelHumidity(90), TempC(45)},
	}

	for _, c := range cases {
		result, err := WetBulbC(c.t, c.rh)
		msgAndArgs := []interface{}{
			"given t %v + rh %v: expected %v, got %v",
			c.t,
			c.rh,
			c.expected,
			result,
		}
		r.NoError(err, msgAndArgs...)
		r.True(eq(result.Unwrap(), c.expected.Unwrap()), msgAndArgs...)
	}
}
