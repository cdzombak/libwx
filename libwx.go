package libwx

// libwx: a variety of weather-related calculations, collected from around the internet
//        and implemented in Go by Chris Dzombak <github.com/cdzombak>.
//
// Also included are some simple type definitions to help avoid erroneously mixing units.

import "math"

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
// and wind speed (in miles/hour). If wind speed is less than 3 mph, or temperature
// if over 50 degrees, the given temperature is returned - the formula works
// below 50 degrees and above 3 mph.
func WindChillF(t TempF, windSpeed SpeedMph) TempF {
	if t > 50.0 || windSpeed < 3.0 {
		return t
	}
	return TempF(35.74 + (0.6215 * float64(t)) - (35.75 * math.Pow(windSpeed.Unwrap(), 0.16)) + (0.4275 * float64(t) * math.Pow(windSpeed.Unwrap(), 0.16)))
}

// WindChillC calculates the wind chill for the given temperature (in Celsius)
// and wind speed (in miles/hour). If wind speed is less than 3 mph, or temperature
// if over 50 degrees, the given temperature is returned - the formula works
// below 50 degrees and above 3 mph.
func WindChillC(temp TempC, windSpeed SpeedMph) TempC {
	return WindChillF(temp.F(), windSpeed).C()
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
