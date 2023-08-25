package libwx

// libwx: a variety of weather-related calculations, collected from around the internet
//        and implemented in Go by Chris Dzombak <github.com/cdzombak>.
//
// Also included are some simple type definitions to help avoid erroneously mixing units.

import "math"

// TempF represents a temperature in degrees Fahrenheit.
type TempF float64

// TempC represents a temperature in degrees Celsius.
type TempC float64

// RelHumidity represents a relative humidity percentage (0-100, inclusive).
type RelHumidity int

// PressureInHg represents barometric pressure in inches of mercury.
type PressureInHg float64

// PressureMb represents barometric pressure in millibars.
type PressureMb float64

// TempFToC converts the given Fahrenheit temperature to Celsius.
func TempFToC(f TempF) TempC {
	return TempC((f - 32.0) / 1.8)
}

// TempCToF converts the given Celsius temperature to Fahrenheit.
func TempCToF(c TempC) TempF {
	return TempF(c*1.8 + 32.0)
}

// DewPointF calculates the dew point given the current temperature (in Fahrenheit)
// and relative humidity percentage (an integer 0-100, *not* a float 0.0-1.0).
func DewPointF(temp TempF, relH RelHumidity) TempF {
	return TempCToF(DewPointC(TempFToC(temp), relH))
}

// DewPointC calculates the dew point given the current temperature (in Celsius)
// and relative humidity percentage (an integer 0-100, *not* a float 0.0-1.0).
func DewPointC(t TempC, relH RelHumidity) TempC {
	const (
		a = 17.625
		b = 243.04
	)
	alpha := math.Log(float64(relH)/100.0) + a*float64(t)/(b+float64(t))
	return TempC((b * alpha) / (a - alpha))
}

// WindChillF calculates the wind chill for the given temperature (in Fahrenheit)
// and wind speed (in miles/hour). If wind speed is less than 3 mph, or temperature
// if over 50 degrees, the given temperature is returned - the formula works
// below 50 degrees and above 3 mph.
func WindChillF(temp TempF, windSpeedMph float64) TempF {
	if temp > 50.0 || windSpeedMph < 3.0 {
		return temp
	}
	return TempF(35.74 + (0.6215 * float64(temp)) - (35.75 * math.Pow(windSpeedMph, 0.16)) + (0.4275 * float64(temp) * math.Pow(windSpeedMph, 0.16)))
}

// WindChillC calculates the wind chill for the given temperature (in Fahrenheit)
// and wind speed (in miles/hour). If wind speed is less than 3 mph, or temperature
// if over 50 degrees, the given temperature is returned - the formula works
// below 50 degrees and above 3 mph.
func WindChillC(temp TempC, windSpeedMph float64) TempC {
	return TempFToC(WindChillF(TempCToF(temp), windSpeedMph))
}

// IndoorHumidityRecommendationF returns the maximum recommended indoor relative
// humidity percentage for the given outdoor temperature (in degrees F).
func IndoorHumidityRecommendationF(outdoorTemp TempF) RelHumidity {
	if outdoorTemp >= 50 {
		return 50
	}
	if outdoorTemp >= 40 {
		return 45
	}
	if outdoorTemp >= 30 {
		return 40
	}
	if outdoorTemp >= 20 {
		return 35
	}
	if outdoorTemp >= 10 {
		return 30
	}
	if outdoorTemp >= 0 {
		return 25
	}
	if outdoorTemp >= -10 {
		return 20
	}
	return 15
}

// IndoorHumidityRecommendationC returns the maximum recommended indoor relative
// humidity percentage for the given outdoor temperature (in degrees C).
func IndoorHumidityRecommendationC(outdoorTemp TempC) RelHumidity {
	return IndoorHumidityRecommendationF(TempCToF(outdoorTemp))
}

// ClampRelHumidity ensures that the given relative humidity percentage is within
// the valid 0-100 inclusive range.
func ClampRelHumidity(h RelHumidity) RelHumidity {
	if h < 0 {
		return 0
	}
	if h > 100 {
		return 100
	}
	return h
}

// PressureMbToInHg converts the given pressure in millibars to inches of mercury.
func PressureMbToInHg(p PressureMb) PressureInHg {
	return PressureInHg(float64(p) / 33.8639)
}

// PressureInHgToMb converts the given pressure in inches of mercury to millibars.
func PressureInHgToMb(p PressureInHg) PressureMb {
	return PressureMb(float64(p) * 33.8639)
}
