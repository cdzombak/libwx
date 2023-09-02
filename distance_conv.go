package libwx

// Km returns the distance in kilometers.
func (mi Mile) Km() Km {
	return Km(mi * 1.60934)
}

// NauticalMiles returns the distance in nautical miles.
func (mi Mile) NauticalMiles() NauticalMile {
	return NauticalMile(mi / 1.15078)
}

// Miles returns the distance in miles.
func (km Km) Miles() Mile {
	return Mile(km / 1.60934)
}

// NauticalMiles returns the distance in nautical miles.
func (km Km) NauticalMiles() NauticalMile {
	return NauticalMile(km / 1.852)
}

// Miles returns the distance in miles.
func (nm NauticalMile) Miles() Mile {
	return Mile(nm * 1.15078)
}

// Km returns the distance in kilometers.
func (nm NauticalMile) Km() Km {
	return Km(nm * 1.852)
}
