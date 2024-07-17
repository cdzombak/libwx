package libwx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO(cdzombak): testing for calculations in libwx.go

func Test_WetBulb_C(t *testing.T) {
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

func Test_HeatIndex_All(t *testing.T) {
	r := require.New(t)
	eq := CurriedFloat64Equal(0.5)

	cases := []struct {
		t        TempF
		rh       RelHumidity
		expected TempF
	}{
		{TempF(80), RelHumidity(40), TempF(80)},
		{TempF(80), RelHumidity(60), TempF(82)},
		{TempF(80), RelHumidity(80), TempF(84)},
		{TempF(80), RelHumidity(100), TempF(87)},
		{TempF(86), RelHumidity(40), TempF(85)},
		{TempF(86), RelHumidity(60), TempF(91)},
		{TempF(86), RelHumidity(80), TempF(100)},
		{TempF(86), RelHumidity(100), TempF(112)},
		{TempF(90), RelHumidity(40), TempF(91)},
		{TempF(90), RelHumidity(60), TempF(100)},
		{TempF(90), RelHumidity(80), TempF(113)},
		{TempF(90), RelHumidity(100), TempF(132)},
		{TempF(104), RelHumidity(40), TempF(119)},
		{TempF(104), RelHumidity(45), TempF(124)},
		{TempF(104), RelHumidity(50), TempF(131)},
		{TempF(104), RelHumidity(55), TempF(137)},
	}

	for _, c := range cases {
		resultF := HeatIndexF(c.t, c.rh)
		msgAndArgsF := []interface{}{
			"given t %v degF + rh %v: expected %v, got %v",
			c.t,
			c.rh,
			c.expected,
			resultF,
		}
		r.True(eq(resultF.Unwrap(), c.expected.Unwrap()), msgAndArgsF...)

		resultC := HeatIndexC(c.t.C(), c.rh)
		msgAndArgsC := []interface{}{
			"given t %v degC + rh %v: expected %v, got %v",
			c.t.C(),
			c.rh,
			c.expected.C(),
			resultC,
		}
		r.True(eq(resultC.Unwrap(), c.expected.C().Unwrap()), msgAndArgsC...)
	}
}

func Test_HeatIndexWarning_F(t *testing.T) {
	r := require.New(t)

	cases := []struct {
		t        TempF
		expected HeatIndexWarning
	}{
		{TempF(-19), HeatIndexWarningNone},
		{TempF(79), HeatIndexWarningNone},
		{TempF(80), HeatIndexWarningCaution},
		{TempF(87), HeatIndexWarningCaution},
		{TempF(90), HeatIndexWarningCaution},
		{TempF(91), HeatIndexWarningExtremeCaution},
		{TempF(103), HeatIndexWarningExtremeCaution},
		{TempF(104), HeatIndexWarningDanger},
		{TempF(124), HeatIndexWarningDanger},
		{TempF(126), HeatIndexWarningExtremeDanger},
		{TempF(135), HeatIndexWarningExtremeDanger},
	}

	for _, c := range cases {
		result := HeatIndexWarningF(c.t)
		msgAndArgs := []interface{}{
			"given t %v degF: expected %v, got %v",
			c.t,
			c.expected,
			result,
		}
		r.Equal(c.expected, result, msgAndArgs...)
	}
}
