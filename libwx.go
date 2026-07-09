package libwx

// libwx: a variety of weather-related calculations, collected from around the internet
//        and implemented in Go by Chris Dzombak <github.com/cdzombak>.
//
// Also included are some simple type definitions to help avoid erroneously mixing units.

import (
	"errors"
	"math"
)

var ErrInputRange = errors.New("one or more input values are outside the calculation's supported range")
var ErrMismatchedInputLength = errors.New("input slices must be the same length")

// DewPointF calculates the dew point given the current temperature (in Fahrenheit)
// and relative humidity percentage (an integer 0-100, *not* a float 0.0-1.0).
func DewPointF(t TempF, rh RelHumidity) TempF {
	return DewPointC(t.C(), rh).F()
}

// DewPointC calculates the dew point given the current temperature (in Celsius)
// and relative humidity percentage (an integer 0-100, *not* a float 0.0-1.0).
func DewPointC(t TempC, rh RelHumidity) TempC {
	rh = rh.Clamped()
	const (
		a = 17.625
		b = 243.04
	)
	alpha := math.Log(float64(rh)/100.0) + a*float64(t)/(b+float64(t))
	return TempC((b * alpha) / (a - alpha))
}

// WindChillF calculates the wind chill for the given temperature (in Fahrenheit)
// and wind speed (in miles/hour).
// If wind speed is less than 3 mph, or temperature is over 50 degrees F, the
// given temperature is returned - the formula works below 50 degrees F and above 3 mph.
func WindChillF(t TempF, windSpeed SpeedMph) TempF {
	if t > 50.0 || windSpeed < 3.0 {
		return t
	}
	return TempF(35.74 + (0.6215 * float64(t)) - (35.75 * math.Pow(windSpeed.Unwrap(), 0.16)) + (0.4275 * float64(t) * math.Pow(windSpeed.Unwrap(), 0.16)))
}

// WindChillFWithValidation calculates the wind chill for the given temperature (in Fahrenheit)
// and wind speed (in miles/hour).
// If wind speed or temperature are outside the supported range, ErrInputRange is returned.
func WindChillFWithValidation(t TempF, windSpeed SpeedMph) (TempF, error) {
	if t > 50.0 || windSpeed < 3.0 {
		return t, ErrInputRange
	}
	return WindChillF(t, windSpeed), nil
}

// WindChillC calculates the wind chill for the given temperature (in Celsius)
// and wind speed (in miles/hour).
// If wind speed is less than 3 mph, or temperature is over 10 degrees C, the
// given temperature is returned - the formula works below 10 degrees C and above 3 mph.
func WindChillC(temp TempC, windSpeed SpeedMph) TempC {
	return WindChillF(temp.F(), windSpeed).C()
}

// WindChillCWithValidation calculates the wind chill for the given temperature (in Celsius)
// and wind speed (in miles/hour).
// If wind speed or temperature are outside the supported range, ErrInputRange is returned.
func WindChillCWithValidation(temp TempC, windSpeed SpeedMph) (TempC, error) {
	if temp.F() > 50.0 || windSpeed < 3.0 {
		return temp, ErrInputRange
	}
	return WindChillF(temp.F(), windSpeed).C(), nil
}

// IndoorHumidityRecommendationF returns the maximum recommended indoor relative
// humidity percentage for the given outdoor temperature (in degrees F).
func IndoorHumidityRecommendationF(outdoorT TempF) RelHumidity {
	if outdoorT >= 50 {
		return 50
	}
	if outdoorT >= 40 {
		return 45
	}
	if outdoorT >= 30 {
		return 40
	}
	if outdoorT >= 20 {
		return 35
	}
	if outdoorT >= 10 {
		return 30
	}
	if outdoorT >= 0 {
		return 25
	}
	if outdoorT >= -10 {
		return 20
	}
	return 15
}

// IndoorHumidityRecommendationC returns the maximum recommended indoor relative
// humidity percentage for the given outdoor temperature (in degrees C).
func IndoorHumidityRecommendationC(outdoorT TempC) RelHumidity {
	return IndoorHumidityRecommendationF(outdoorT.F())
}

// WetBulbF calculates the wet bulb temperature (in Fahrenheit) given a dry bulb
// temperature (in Fahrenheit) and relative humidity percentage.
// If the given temperature or relative humidity are outside the supported range,
// ErrInputRange is returned.
// See: https://journals.ametsoc.org/view/journals/apme/50/11/jamc-d-11-0143.1.xml
func WetBulbF(temp TempF, rh RelHumidity) (TempF, error) {
	result, err := WetBulbC(temp.C(), rh)
	return result.F(), err
}

// WetBulbC calculates the wet bulb temperature (in Celsius) given a dry bulb
// temperature (in Celsius) and relative humidity percentage.
// If the given temperature or relative humidity are outside the supported range,
// ErrInputRange is returned.
// See: https://journals.ametsoc.org/view/journals/apme/50/11/jamc-d-11-0143.1.xml
func WetBulbC(temp TempC, rh RelHumidity) (TempC, error) {
	rh = rh.Clamped()
	if rh.Unwrap() < 5 || rh.Unwrap() > 99 {
		return temp, ErrInputRange
	}
	if temp.Unwrap() < -20 || temp.Unwrap() > 50 {
		return temp, ErrInputRange
	}
	// formula for the left validity border line:
	// y = -1*(75-5)/(20+9.25) * x + 25
	// where y == RH% and x == dry bulb degC
	// see full-jamc-d-11-0143.1-f3-annotated.jpg in this repo, taken from
	// https://journals.ametsoc.org/view/journals/apme/50/11/jamc-d-11-0143.1.xml
	// and annotated
	y := -1*(75-5)/(20+9.25)*temp.Unwrap() + 25
	if y > rh.UnwrapFloat64() {
		return temp, ErrInputRange
	}

	// Tw = T*atan[0.151977(RH% + 8.313659)**1/2] + atan(T + RH%) - atan(RH% - 1.676331)
	// + 0.00391838(RH%)**3/2 * atan(0.023101*RH%) - 4.686035
	// taken from figure 1 of https://journals.ametsoc.org/view/journals/apme/50/11/jamc-d-11-0143.1.xml
	return TempC(
			temp.Unwrap()*math.Atan(0.151977*math.Pow(rh.UnwrapFloat64()+8.313659, 0.5)) +
				math.Atan(temp.Unwrap()+rh.UnwrapFloat64()) -
				math.Atan(rh.UnwrapFloat64()-1.676331) +
				0.00391838*math.Pow(rh.UnwrapFloat64(), 1.5)*math.Atan(0.023101*rh.UnwrapFloat64()) -
				4.686035,
		),
		nil
}

const (
	// MinWetBulbPressureMb is the minimum barometric pressure supported by
	// WetBulbFAtPressure and WetBulbCAtPressure.
	MinWetBulbPressureMb = PressureMb(585.0)

	// MaxWetBulbPressureMb is the maximum barometric pressure supported by
	// WetBulbFAtPressure and WetBulbCAtPressure.
	MaxWetBulbPressureMb = PressureMb(1100.0)
)

// WetBulbFAtPressure calculates the wet bulb temperature (in Fahrenheit) given a
// dry bulb temperature (in Fahrenheit), relative humidity percentage, and
// barometric pressure (in millibars).
// If the given temperature, relative humidity, or pressure are outside the
// supported range, ErrInputRange is returned.
// See: https://journals.ametsoc.org/view/journals/atot/30/8/jtech-d-12-00191_1.xml
func WetBulbFAtPressure(temp TempF, rh RelHumidity, pressure PressureMb) (TempF, error) {
	result, err := WetBulbCAtPressure(temp.C(), rh, pressure)
	return result.F(), err
}

// WetBulbCAtPressure calculates the wet bulb temperature (in Celsius) given a
// dry bulb temperature (in Celsius), relative humidity percentage, and
// barometric pressure (in millibars).
// The supported input ranges are -17 to 40 degrees C, 5-99% relative humidity,
// and MinWetBulbPressureMb to MaxWetBulbPressureMb (inclusive); if the input is
// outside these ranges, ErrInputRange is returned.
// The formula is taken from "Direct Calculation of Thermodynamic Wet-Bulb
// Temperature as a Function of Pressure and Elevation" (Sadeghi et al., Journal
// of Atmospheric and Oceanic Technology, 2013). Its stated accuracy is within
// 0.65 degrees C, with the largest errors occurring at high temperatures
// combined with very low or very high relative humidity.
// See: https://journals.ametsoc.org/view/journals/atot/30/8/jtech-d-12-00191_1.xml
func WetBulbCAtPressure(temp TempC, rh RelHumidity, pressure PressureMb) (TempC, error) {
	rh = rh.Clamped()
	if rh.Unwrap() < 5 || rh.Unwrap() > 99 {
		return temp, ErrInputRange
	}
	if temp.Unwrap() < -17 || temp.Unwrap() > 40 {
		return temp, ErrInputRange
	}
	if pressure < MinWetBulbPressureMb || pressure > MaxWetBulbPressureMb {
		return temp, ErrInputRange
	}

	// equation references below are to Sadeghi et al. 2013,
	// https://journals.ametsoc.org/view/journals/atot/30/8/jtech-d-12-00191_1.xml
	// named as in the paper: Ta is the dry bulb temperature, Pa the pressure.
	// the paper's equations work in kPa; 1 millibar == 0.1 kPa
	ta := temp.Unwrap()
	pa := pressure.Unwrap() / 10.0

	// thermodynamic psychrometric constant (degC**-1).
	// the paper defines gamma = Cp/Lv, approximately 0.4 g_water/(kg_air * K);
	// dividing by the 0.622 ratio of the molecular weights of water and dry air
	// expresses it in terms of vapor pressure: 0.4e-3 / 0.622 == 6.43e-4.
	// this agrees with the upper end of the 5.68e-4 to 6.42e-4 range of empirical
	// estimates the paper cites, and yields gamma*Pa == 0.0650 kPa/degC at standard
	// sea level pressure, close to the standard psychrometric constant of 0.0665.
	const gamma = 6.42e-4

	// saturated vapor pressure (kPa) at the dry bulb temperature, via the
	// Magnus equation (equation 5a) with Buck 1981 coefficients (table 1).
	// over the supported -17 to 40 degC range the two coefficient sets differ
	// by at most 0.011 degC in the final result, so this branch cannot be
	// meaningfully covered by a test; it is kept for fidelity to table 1.
	b, c := 17.368, 238.88 // valid 0 to 50 degC
	if ta < 0 {
		b, c = 17.966, 247.15 // valid -40 to 0 degC
	}
	esTa := 0.611 * math.Exp(b*ta/(ta+c))
	// ambient vapor pressure (kPa)
	ea := esTa * rh.UnwrapFloat64() / 100.0

	// lambda = 0.0014*exp(0.027*Ta) (equation 11)
	lambda := 0.0014 * math.Exp(0.027*ta)
	// zeta = -3e-7*Ta**3 - 1e-5*Ta**2 + 2e-5*Ta + 4.44e-2 (equation 12)
	zeta := -3e-7*ta*ta*ta - 1e-5*ta*ta + 2e-5*ta + 4.44e-2
	// phi = zeta + gamma*Pa (equation 9c)
	phi := zeta + gamma*pa
	// psi = a - gamma*Pa*Ta - ea, where a == 0.611 kPa (equation 9a)
	psi := 0.611 - gamma*pa*ta - ea

	// Tw = (-phi + sqrt(phi**2 - 4*lambda*psi)) / (2*lambda) (equations 8a/8b);
	// the positive root is always the correct one (see discussion after equation 8b)
	wetBulb := (-phi + math.Sqrt(phi*phi-4*lambda*psi)) / (2 * lambda)

	// regression correction mapping the quadratic solution to the theoretical
	// wet bulb temperature, from figure 3
	if wetBulb < 0 {
		wetBulb = 1.011*wetBulb + 0.0419
	} else {
		wetBulb = 1.0301*wetBulb - 0.213
	}

	// near saturation the figure 3 correction can push the result slightly
	// above the dry bulb temperature, which is physically impossible
	return TempC(math.Min(wetBulb, ta)), nil
}

func heatIndexConstantsF() [9]float64 {
	// from https://en.wikipedia.org/wiki/Heat_index#Formula
	// captured on 2024-07-17
	return [9]float64{
		-42.379,
		2.04901523,
		10.14333127,
		-0.22475541,
		-6.83783e-3,
		-5.481717e-2,
		1.22874e-3,
		8.5282e-4,
		-1.99e-6,
	}
}

func heatIndexConstantsC() [9]float64 {
	// from https://en.wikipedia.org/wiki/Heat_index#Formula
	// captured on 2024-07-17
	return [9]float64{
		-8.78469475556,
		1.61139411,
		2.33854883889,
		-0.14611605,
		-0.012308094,
		-0.0164248277778,
		2.211732e-3,
		7.2546e-4,
		-3.582e-6,
	}
}

func rawHeatIndex(c [9]float64, rawTemp, rawRelH float64) float64 {
	// from https://en.wikipedia.org/wiki/Heat_index#Formula
	// captured on 2024-07-17
	// note that constants on that page are 1-indexed, while the constants
	// array here is 0-indexed
	// see also: https://www.weather.gov/media/ffc/ta_htindx.PDF
	return c[0] +
		c[1]*rawTemp +
		c[2]*rawRelH +
		c[3]*rawTemp*rawRelH +
		c[4]*math.Pow(rawTemp, 2) +
		c[5]*math.Pow(rawRelH, 2) +
		c[6]*math.Pow(rawTemp, 2)*rawRelH +
		c[7]*rawTemp*math.Pow(rawRelH, 2) +
		c[8]*math.Pow(rawTemp, 2)*math.Pow(rawRelH, 2)
}

// HeatIndexF is deprecated; use HeatIndexFWithValidation
func HeatIndexF(temp TempF, rh RelHumidity) TempF {
	retv, _ := HeatIndexFWithValidation(temp, rh)
	return retv
}

// HeatIndexC is deprecated; use HeatIndexCWithValidation
func HeatIndexC(temp TempC, rh RelHumidity) TempC {
	retv, _ := HeatIndexCWithValidation(temp, rh)
	return retv
}

// HeatIndexFWithValidation calculates the heat index for the given temperature (in Fahrenheit)
// and relative humidity percentage.
func HeatIndexFWithValidation(temp TempF, rh RelHumidity) (TempF, error) {
	var err error
	if temp < TempC(25).F() {
		err = ErrInputRange
	}
	return TempF(rawHeatIndex(
		heatIndexConstantsF(),
		temp.Unwrap(),
		rh.Clamped().UnwrapFloat64(),
	)), err
}

// HeatIndexCWithValidation calculates the heat index for the given temperature (in Celsius)
// and relative humidity percentage.
func HeatIndexCWithValidation(temp TempC, rh RelHumidity) (TempC, error) {
	var err error
	if temp < TempC(25) {
		err = ErrInputRange
	}
	return TempC(rawHeatIndex(
		heatIndexConstantsC(),
		temp.Unwrap(),
		rh.Clamped().UnwrapFloat64(),
	)), err
}

// HeatIndexWarningF returns a heat index warning level for the
// given heat index temperature (in Fahrenheit) per
// https://en.wikipedia.org/wiki/Heat_index#Table_of_values
// captured on 2024-07-17.
func HeatIndexWarningF(heatIndex TempF) HeatIndexWarning {
	if heatIndex.Unwrap() < 80 {
		return HeatIndexWarningNone
	}
	if heatIndex.Unwrap() < 91 {
		return HeatIndexWarningCaution
	}
	if heatIndex.Unwrap() < 104 {
		return HeatIndexWarningExtremeCaution
	}
	if heatIndex.Unwrap() < 125 {
		return HeatIndexWarningDanger
	}
	return HeatIndexWarningExtremeDanger
}

// HeatIndexWarningC returns a heat index warning level for the
// given heat index temperature (in Celsius) per
// https://en.wikipedia.org/wiki/Heat_index#Table_of_values
// captured on 2024-07-17.
func HeatIndexWarningC(heatIndex TempC) HeatIndexWarning {
	if heatIndex.Unwrap() < 27 {
		return HeatIndexWarningNone
	}
	if heatIndex.Unwrap() < 33 {
		return HeatIndexWarningCaution
	}
	if heatIndex.Unwrap() < 40 {
		return HeatIndexWarningExtremeCaution
	}
	if heatIndex.Unwrap() < 52 {
		return HeatIndexWarningDanger
	}
	return HeatIndexWarningExtremeDanger
}

// AvgDirectionDeg calculates the circular mean of the given set of angles (in degrees).
// This is useful to find e.g. the average wind direction.
func AvgDirectionDeg(degrees []Degree) Degree {
	return radToDeg(circularMean(degToRadSlice(clampedDegSlice(degrees)), nil)).Clamped()
}

// WeightedAvgDirectionDeg calculates the weighted circular mean of the given set of angles (in degrees).
// This is useful to find e.g. the average wind direction, weighted by wind speed.
func WeightedAvgDirectionDeg(degrees []Degree, weights []float64) (Degree, error) {
	if len(degrees) != len(weights) {
		return 0.0, ErrMismatchedInputLength
	}
	return radToDeg(circularMean(degToRadSlice(clampedDegSlice(degrees)), weights)).Clamped(), nil
}

// StdDevDirectionDeg calculates the circular standard deviation of the given set of angles (in degrees).
// This is useful to find e.g. the variability of wind direction.
func StdDevDirectionDeg(degrees []Degree) Degree {
	return radToDeg(circularStdDev(degToRadSlice(clampedDegSlice(degrees)), nil))
}

// WeightedStdDevDirectionDeg calculates the circular standard deviation of the given set of angles (in degrees).
// This is useful to find e.g. the variability of wind direction, weighted by wind speed.
func WeightedStdDevDirectionDeg(degrees []Degree, weights []float64) (Degree, error) {
	if len(degrees) != len(weights) {
		return 0.0, ErrMismatchedInputLength
	}
	return radToDeg(circularStdDev(degToRadSlice(clampedDegSlice(degrees)), weights)), nil
}

func saturationVaporPressureC(temp TempC) float64 {
	if temp.Unwrap() < -20 || temp.Unwrap() > 100 {
		return 0
	}
	const (
		A = 8.07131
		B = 1730.63
		C = 233.426
	)
	logP := A - B/(C+temp.Unwrap())
	return math.Pow(10, logP)
}

func AbsHumidityFromRelF(temp TempF, rh RelHumidity) AbsHumidity {
	return AbsHumidityFromRelC(temp.C(), rh)
}

func AbsHumidityFromRelC(temp TempC, rh RelHumidity) AbsHumidity {
	rh = rh.Clamped()

	pSat := saturationVaporPressureC(temp)
	if pSat == 0 {
		return 0
	}

	pSatPa := pSat * 133.322

	tempK := temp.Unwrap() + 273.15
	ah := (rh.UnwrapFloat64() / 100.0) * (pSatPa * 18.016) / (8.314 * tempK)

	return AbsHumidity(ah)
}

func RelHumidityFromAbsF(temp TempF, ah AbsHumidity) RelHumidity {
	return RelHumidityFromAbsC(temp.C(), ah)
}

func RelHumidityFromAbsC(temp TempC, ah AbsHumidity) RelHumidity {
	pSat := saturationVaporPressureC(temp)
	if pSat == 0 {
		return 0
	}

	pSatPa := pSat * 133.322

	tempK := temp.Unwrap() + 273.15
	rh := (ah.Unwrap() * 8.314 * tempK) / (pSatPa * 18.016) * 100.0

	return ClampedRelHumidity(int(rh + 0.5))
}
