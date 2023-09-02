package libwx

// PressureInHg represents barometric pressure in inches of mercury.
type PressureInHg float64

// PressureMb represents barometric pressure in millibars.
type PressureMb float64

func (p PressureMb) Unwrap() float64   { return float64(p) }
func (p PressureInHg) Unwrap() float64 { return float64(p) }
