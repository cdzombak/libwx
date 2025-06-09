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
	ah = ah.Clamped()

	pSat := saturationVaporPressureC(temp)

	if pSat == 0 {
		return 0
	}

	pSatPa := pSat * 133.322

	tempK := temp.Unwrap() + 273.15
	rh := (ah.Unwrap() * 8.314 * tempK) / (pSatPa * 18.016) * 100.0

	return ClampedRelHumidity(int(rh + 0.5))
}
