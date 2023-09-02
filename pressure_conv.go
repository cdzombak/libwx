package libwx

// InHg converts pressure in millibars to inches of mercury.
func (p PressureMb) InHg() PressureInHg {
	return PressureInHg(p / 33.8639)
}

// Mb converts pressure in inches of mercury to millibars.
func (p PressureInHg) Mb() PressureMb {
	return PressureMb(p * 33.8639)
}
