package libwx

// SpeedMph represents speed (e.g. wind speed) in miles per hour.
type SpeedMph float64

// SpeedKmH represents speed in kilometers per hour.
type SpeedKmH float64

// SpeedKnots represents speed in knots.
type SpeedKnots float64

func (s SpeedMph) Unwrap() float64   { return float64(s) }
func (s SpeedKmH) Unwrap() float64   { return float64(s) }
func (s SpeedKnots) Unwrap() float64 { return float64(s) }
