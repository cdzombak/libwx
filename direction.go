package libwx

import "math"

type DirectionStrPrecision int

const (
	DirectionStrPrecision1 DirectionStrPrecision = 1
	DirectionStrPrecision2 DirectionStrPrecision = 2
	DirectionStrPrecision3 DirectionStrPrecision = 3
)

// DirectionStr returns a string representation of the given compass direction (in degrees).
// Given DirectionStrPrecision1, it returns a cardinal direction (N, E, S, W).
// Given DirectionStrPrecision2, it returns a primary intercardinal direction (N, NE, E, SE, S, SW, W, NW).
// Given DirectionStrPrecision3, it returns a secondary intercardinal direction (N, NNE, NE, ENE, E, ESE, SE, SSE, S, SSW, SW, WSW, W, WNW, NW, NNW).
func DirectionStr(deg Degree, precision DirectionStrPrecision) string {
	deg = deg.Clamped()

	switch precision {
	case DirectionStrPrecision3:
		return []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}[int((deg/22.5)+.5)%16]
	case DirectionStrPrecision2:
		return []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}[int((deg/45.0)+.5)%8]
	case DirectionStrPrecision1:
		fallthrough
	default:
		return []string{"N", "E", "S", "W"}[int((deg/90.0)+.5)%4]
	}
}

// ClampedDegree returns a Degree from the given value, guaranteed
// to be within 0 <= d < 360.
func ClampedDegree(d float64) Degree {
	return Degree(d).Clamped()
}

// Clamped returns an angular direction in degrees guaranteed to be within 0 <= d < 360.
func (d Degree) Clamped() Degree {
	for d < 0 {
		d += 360
	}
	for d >= 360 {
		d -= 360
	}
	return d
}

func clampedDegSlice(in []Degree) []Degree {
	retv := make([]Degree, len(in))
	for i, v := range in {
		retv[i] = v.Clamped()
	}
	return retv
}

func degToRadSlice(in []Degree) []float64 {
	retv := make([]float64, len(in))
	for i, v := range in {
		retv[i] = degToRad(v)
	}
	return retv
}

func degToRad(deg Degree) float64 {
	return deg.Unwrap() * math.Pi / 180.0
}

func radToDeg(rad float64) Degree {
	return Degree(rad * 180.0 / math.Pi)
}
