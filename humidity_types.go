package libwx

// RelHumidity represents a relative humidity percentage (0-100, inclusive).
type RelHumidity int

func (rh RelHumidity) Unwrap() int            { return int(rh) }
func (rh RelHumidity) UnwrapFloat64() float64 { return float64(rh) }

// ClampedRelHumidity returns a RelHumidity from the given integer, guaranteed
// to be within the valid 0-100 (inclusive) range.
func ClampedRelHumidity(rh int) RelHumidity {
	return RelHumidity(rh).Clamped()
}

// Clamped returns a relative humidity guaranteed to be within
// the valid 0-100 (inclusive) range.
func (rh RelHumidity) Clamped() RelHumidity {
	if rh < 0 {
		return 0
	}
	if rh > 100 {
		return 100
	}
	return rh
}

type AbsHumidity float64

func (ah AbsHumidity) Unwrap() float64 { return float64(ah) }
