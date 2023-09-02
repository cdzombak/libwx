package libwx

// Mile represents distance in miles.
type Mile float64

// Km represents distance in kilometers.
type Km float64

// NauticalMile represents distance in nautical miles.
type NauticalMile float64

func (mi Mile) Unwrap() float64         { return float64(mi) }
func (km Km) Unwrap() float64           { return float64(km) }
func (nm NauticalMile) Unwrap() float64 { return float64(nm) }
