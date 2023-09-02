package libwx

// KmH returns the speed in kilometers per hour.
func (s SpeedMph) KmH() SpeedKmH {
	return SpeedKmH(s * 1.60934)
}

// Knots returns the speed in knots.
func (s SpeedMph) Knots() SpeedKnots {
	return SpeedKnots(s / 1.15078)
}

// Mph returns the speed in miles per hour.
func (s SpeedKmH) Mph() SpeedMph {
	return SpeedMph(s / 1.60934)
}

// Knots returns the speed in knots.
func (s SpeedKmH) Knots() SpeedKnots {
	return SpeedKnots(s / 1.852)
}

// Mph returns the speed in miles per hour.
func (s SpeedKnots) Mph() SpeedMph {
	return SpeedMph(s * 1.15078)
}

// KmH returns the speed in kilometers per hour.
func (s SpeedKnots) KmH() SpeedKmH {
	return SpeedKmH(s * 1.852)
}
