package libwx

import (
	"math"
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

// wetBulbTheoreticalC solves the psychrometric equation (equation 4 of Sadeghi
// et al. 2013) for the wet bulb temperature by bisection, using the same
// saturated vapor pressure formula and psychrometric constant as
// WetBulbCAtPressure. This is the "theoretical" solution that the paper's
// direct formula approximates; it is used only to verify that approximation.
func wetBulbTheoreticalC(ta float64, rh float64, pressureMb float64) float64 {
	const gamma = 6.42e-4
	es := func(t float64) float64 {
		b, c := 17.368, 238.88
		if t < 0 {
			b, c = 17.966, 247.15
		}
		return 0.611 * math.Exp(b*t/(t+c))
	}
	pa := pressureMb / 10.0
	ea := es(ta) * rh / 100.0
	// ea - es(Tw) + gamma*Pa*(Ta - Tw) is strictly decreasing in Tw, and is
	// non-negative at Tw == Ta, so the root lies in [-80, Ta]
	lo, hi := -80.0, ta
	for i := 0; i < 200; i++ {
		mid := (lo + hi) / 2
		if ea-es(mid)+gamma*pa*(ta-mid) > 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

func Test_WetBulbAtPressure_C(t *testing.T) {
	r := require.New(t)
	eq := CurriedFloat64Equal(Tolerance001)

	// these are the values WetBulbCAtPressure is expected to produce, asserted
	// tightly so that any change to the formula's coefficients is caught. That
	// the formula is *correct* - that these values approximate the theoretical
	// wet bulb temperature - is verified separately, by
	// Test_WetBulbAtPressure_C_MatchesTheoreticalSolution.
	cases := []struct {
		t        TempC
		rh       RelHumidity
		p        PressureMb
		expected TempC
	}{
		// standard sea level pressure, across the supported temperature range
		{TempC(-15), RelHumidity(80), PressureMb(1013.25), TempC(-15.4413)},
		{TempC(-10), RelHumidity(50), PressureMb(1013.25), TempC(-11.7208)},
		{TempC(-5), RelHumidity(90), PressureMb(1013.25), TempC(-5.4503)},
		{TempC(0), RelHumidity(50), PressureMb(1013.25), TempC(-2.8887)},
		{TempC(5), RelHumidity(60), PressureMb(1013.25), TempC(1.9603)},
		{TempC(10), RelHumidity(40), PressureMb(1013.25), TempC(4.4710)},
		{TempC(15), RelHumidity(70), PressureMb(1013.25), TempC(12.0039)},
		{TempC(20), RelHumidity(50), PressureMb(1013.25), TempC(13.8423)},
		{TempC(25), RelHumidity(30), PressureMb(1013.25), TempC(14.4432)},
		{TempC(30), RelHumidity(70), PressureMb(1013.25), TempC(25.6201)},
		{TempC(35), RelHumidity(50), PressureMb(1013.25), TempC(25.9263)},
		{TempC(40), RelHumidity(90), PressureMb(1013.25), TempC(38.4670)},
		// relative humidity boundary values; at 99% the Tw <= Ta clamp binds
		{TempC(20), RelHumidity(5), PressureMb(1013.25), TempC(6.7337)},
		{TempC(20), RelHumidity(99), PressureMb(1013.25), TempC(20.0)},
		// reduced pressures (roughly 1000/2000/3000/4500 m elevation)
		{TempC(20), RelHumidity(50), PressureMb(899), TempC(13.5440)},
		{TempC(20), RelHumidity(50), PressureMb(795), TempC(13.2451)},
		{TempC(20), RelHumidity(50), PressureMb(703), TempC(12.9554)},
		{TempC(20), RelHumidity(50), PressureMb(585), TempC(12.5426)},
		{TempC(0), RelHumidity(70), PressureMb(703), TempC(-2.0988)},
		{TempC(-10), RelHumidity(60), PressureMb(795), TempC(-11.6377)},
		{TempC(30), RelHumidity(40), PressureMb(899), TempC(19.6323)},
		{TempC(5), RelHumidity(90), PressureMb(585), TempC(4.0071)},
		// above-standard sea level pressure
		{TempC(15), RelHumidity(60), PressureMb(1040), TempC(10.9121)},
		{TempC(25), RelHumidity(80), PressureMb(1040), TempC(22.6573)},
		// pressure range bounds are inclusive
		{TempC(20), RelHumidity(50), MaxWetBulbPressureMb, TempC(14.0506)},
		{TempC(20), RelHumidity(50), MinWetBulbPressureMb, TempC(12.5426)},
	}

	for _, c := range cases {
		result, err := WetBulbCAtPressure(c.t, c.rh, c.p)
		msgAndArgs := []interface{}{
			"given t %v + rh %v + p %v: expected %v, got %v",
			c.t,
			c.rh,
			c.p,
			c.expected,
			result,
		}
		r.NoError(err, msgAndArgs...)
		r.True(eq(result.Unwrap(), c.expected.Unwrap()), msgAndArgs...)
		r.LessOrEqual(result.Unwrap(), c.t.Unwrap(), msgAndArgs...)
	}
}

func Test_WetBulbAtPressure_C_MatchesTheoreticalSolution(t *testing.T) {
	r := require.New(t)

	// Sadeghi et al. 2013 claim their direct formula stays within 0.65 degC of
	// the theoretical wet bulb temperature. Across the whole domain we accept,
	// the largest deviation is 0.6703 degC, at the physically implausible corner
	// of 40 degC at 585mb (~4500m elevation); the paper itself notes its largest
	// errors occur at high temperature combined with very low relative humidity.
	// Assert the slightly looser bound that actually holds everywhere.
	const maxDeviationC = 0.7

	for temp := -17; temp <= 40; temp++ {
		for _, p := range []PressureMb{585, 703, 795, 899, 1013.25, 1100} {
			for rh := 5; rh <= 99; rh++ {
				result, err := WetBulbCAtPressure(TempC(temp), RelHumidity(rh), p)
				r.NoError(err)
				theoretical := wetBulbTheoreticalC(float64(temp), float64(rh), p.Unwrap())
				r.InDelta(
					theoretical, result.Unwrap(), maxDeviationC,
					"given t %v + rh %v + p %v: theoretical %v, got %v",
					temp, rh, p, theoretical, result,
				)
			}
		}
	}
}

func Test_WetBulbAtPressure_C_SeaLevelMatchesStull(t *testing.T) {
	r := require.New(t)
	eq := CurriedFloat64Equal(Tolerance0)

	// the grid stops at 35 degC: Sadeghi et al. 2013 report their largest
	// disagreement with the Stull formula at >= 39 degC, where the combined
	// error of the two approximations exceeds 1 degC
	for temp := -15; temp <= 35; temp += 5 {
		for rh := 30; rh <= 95; rh += 5 {
			stull, err := WetBulbC(TempC(temp), RelHumidity(rh))
			if err != nil {
				continue // outside the Stull formula's validity region
			}
			result, err := WetBulbCAtPressure(TempC(temp), RelHumidity(rh), PressureMb(1013.25))
			msgAndArgs := []interface{}{
				"given t %v + rh %v: WetBulbC %v, WetBulbCAtPressure %v",
				temp,
				rh,
				stull,
				result,
			}
			r.NoError(err, msgAndArgs...)
			r.True(eq(result.Unwrap(), stull.Unwrap()), msgAndArgs...)
		}
	}
}

func Test_WetBulbAtPressure_C_PressureMonotonic(t *testing.T) {
	r := require.New(t)

	const (
		dryBulb = TempC(25)
		rh      = RelHumidity(60)
		steps   = 100
	)
	minP, maxP := MinWetBulbPressureMb.Unwrap(), MaxWetBulbPressureMb.Unwrap()

	prev, err := WetBulbCAtPressure(dryBulb, rh, MaxWetBulbPressureMb)
	r.NoError(err)

	// walk the full supported pressure range, ending exactly on
	// MinWetBulbPressureMb so that both bounds are exercised
	var p PressureMb
	for i := 1; i <= steps; i++ {
		p = PressureMb(maxP - (maxP-minP)*float64(i)/float64(steps))
		result, err := WetBulbCAtPressure(dryBulb, rh, p)
		r.NoError(err)
		r.LessOrEqual(result.Unwrap(), prev.Unwrap(), "wet bulb temperature must not increase as pressure decreases (at %v mb)", p)
		r.LessOrEqual(result.Unwrap(), dryBulb.Unwrap(), "wet bulb temperature must not exceed dry bulb temperature (at %v mb)", p)
		prev = result
	}
	r.Equal(MinWetBulbPressureMb, p, "loop must finish at the minimum supported pressure")
}

func Test_WetBulbAtPressure_F(t *testing.T) {
	r := require.New(t)
	eq := CurriedFloat64Equal(Tolerance001)

	cases := []struct {
		t  TempF
		rh RelHumidity
		p  PressureMb
	}{
		{TempF(68), RelHumidity(50), PressureMb(1013.25)},
		{TempF(14), RelHumidity(60), PressureMb(795)},
		{TempF(95), RelHumidity(80), PressureMb(899)},
	}

	for _, c := range cases {
		resultF, err := WetBulbFAtPressure(c.t, c.rh, c.p)
		r.NoError(err)
		resultC, err := WetBulbCAtPressure(c.t.C(), c.rh, c.p)
		r.NoError(err)
		r.True(
			eq(resultF.Unwrap(), resultC.F().Unwrap()),
			"given t %v + rh %v + p %v: expected %v, got %v",
			c.t,
			c.rh,
			c.p,
			resultC.F(),
			resultF,
		)
	}
}

func Test_WetBulbAtPressure_InputRangeErrors(t *testing.T) {
	r := require.New(t)

	cases := []struct {
		t  TempC
		rh RelHumidity
		p  PressureMb
	}{
		// relative humidity outside 5-99%
		{TempC(20), RelHumidity(4), PressureMb(1013.25)},
		{TempC(20), RelHumidity(100), PressureMb(1013.25)},
		// temperature outside -17 to 40 degC
		{TempC(-18), RelHumidity(50), PressureMb(1013.25)},
		{TempC(41), RelHumidity(50), PressureMb(1013.25)},
		// pressure outside MinWetBulbPressureMb to MaxWetBulbPressureMb
		{TempC(20), RelHumidity(50), MinWetBulbPressureMb - 1},
		{TempC(20), RelHumidity(50), MaxWetBulbPressureMb + 1},
	}

	for _, c := range cases {
		result, err := WetBulbCAtPressure(c.t, c.rh, c.p)
		msgAndArgs := []interface{}{
			"given t %v + rh %v + p %v",
			c.t,
			c.rh,
			c.p,
		}
		r.ErrorIs(err, ErrInputRange, msgAndArgs...)
		r.Equal(c.t, result, msgAndArgs...)

		resultF, err := WetBulbFAtPressure(c.t.F(), c.rh, c.p)
		r.ErrorIs(err, ErrInputRange, msgAndArgs...)
		r.InDelta(c.t.F().Unwrap(), resultF.Unwrap(), Tolerance001, msgAndArgs...)
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

func Test_AbsHumidity_Conversions(t *testing.T) {
	r := require.New(t)
	eq := CurriedFloat64Equal(0.5)

	cases := []struct {
		tempC    TempC
		rh       RelHumidity
		expected AbsHumidity
		desc     string
	}{
		{TempC(20), RelHumidity(50), AbsHumidity(8.7), "20°C, 50% RH - typical indoor conditions"},
		{TempC(25), RelHumidity(60), AbsHumidity(13.8), "25°C, 60% RH - warm indoor conditions"},
		{TempC(30), RelHumidity(70), AbsHumidity(21.1), "30°C, 70% RH - hot humid conditions"},
		{TempC(10), RelHumidity(80), AbsHumidity(7.6), "10°C, 80% RH - cool humid conditions"},
		{TempC(0), RelHumidity(90), AbsHumidity(4.3), "0°C, 90% RH - near freezing"},
		{TempC(-10), RelHumidity(95), AbsHumidity(2.2), "-10°C, 95% RH - cold conditions"},
	}

	for _, c := range cases {
		result := AbsHumidityFromRelC(c.tempC, c.rh)
		msgAndArgs := []interface{}{
			"AbsHumidityFromRelC: %s - given t %v + rh %v: expected %v g/m³, got %v g/m³",
			c.desc, c.tempC, c.rh, c.expected, result,
		}
		r.True(eq(result.Unwrap(), c.expected.Unwrap()), msgAndArgs...)

		resultRh := RelHumidityFromAbsC(c.tempC, c.expected)
		msgAndArgsRh := []interface{}{
			"RelHumidityFromAbsC: %s - given t %v + ah %v: expected %v%%, got %v%%",
			c.desc, c.tempC, c.expected, c.rh, resultRh,
		}
		r.True(IntCompare(resultRh.Unwrap(), c.rh.Unwrap()) == 0 ||
			IntCompare(resultRh.Unwrap(), c.rh.Unwrap()-1) == 0 ||
			IntCompare(resultRh.Unwrap(), c.rh.Unwrap()+1) == 0, msgAndArgsRh...)
	}
}

func Test_AbsHumidity_Fahrenheit(t *testing.T) {
	r := require.New(t)
	eq := CurriedFloat64Equal(0.5)

	cases := []struct {
		tempF TempF
		tempC TempC
		rh    RelHumidity
	}{
		{TempF(68), TempC(20), RelHumidity(50)},
		{TempF(77), TempC(25), RelHumidity(60)},
		{TempF(32), TempC(0), RelHumidity(90)},
	}

	for _, c := range cases {
		resultF := AbsHumidityFromRelF(c.tempF, c.rh)
		resultC := AbsHumidityFromRelC(c.tempC, c.rh)

		msgAndArgs := []interface{}{
			"Fahrenheit and Celsius should give same result: F=%v, C=%v, RH=%v - got F=%v, C=%v",
			c.tempF, c.tempC, c.rh, resultF, resultC,
		}
		r.True(eq(resultF.Unwrap(), resultC.Unwrap()), msgAndArgs...)
	}
}

func Test_AbsHumidity_EdgeCases(t *testing.T) {
	r := require.New(t)

	result := AbsHumidityFromRelC(TempC(-30), RelHumidity(50))
	r.Equal(AbsHumidity(0), result, "Should return 0 for temperature out of range")

	result = AbsHumidityFromRelC(TempC(110), RelHumidity(50))
	r.Equal(AbsHumidity(0), result, "Should return 0 for temperature out of range")
}
