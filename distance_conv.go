package libwx

// Km returns the distance in kilometers.
func (mi Mile) Km() Km {
	return Km(mi * 1.60934)
}

// NauticalMiles returns the distance in nautical miles.
func (mi Mile) NauticalMiles() NauticalMile {
	return NauticalMile(mi / 1.15078)
}

// Meters returns the distance in meters.
func (mi Mile) Meters() Meter {
	return Meter(mi * 1609.34)
}

// Miles returns the distance in miles.
func (km Km) Miles() Mile {
	return Mile(km / 1.60934)
}

// NauticalMiles returns the distance in nautical miles.
func (km Km) NauticalMiles() NauticalMile {
	return NauticalMile(km / 1.852)
}

// Meters returns the distance in meters.
func (km Km) Meters() Meter {
	return Meter(km * 1000)
}

// Miles returns the distance in miles.
func (nm NauticalMile) Miles() Mile {
	return Mile(nm * 1.15078)
}

// Km returns the distance in kilometers.
func (nm NauticalMile) Km() Km {
	return Km(nm * 1.852)
}

// Meters returns the distance in meters.
func (nm NauticalMile) Meters() Meter {
	return Meter(nm * 1852)
}

// Miles returns the distance in miles.
func (m Meter) Miles() Mile {
	return Mile(m / 1609.34)
}

// Km returns the distance in kilometers.
func (m Meter) Km() Km {
	return Km(m / 1000)
}

// NauticalMiles returns the distance in nautical miles.
func (m Meter) NauticalMiles() NauticalMile {
	return NauticalMile(m / 1852)
}
