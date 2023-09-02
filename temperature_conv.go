package libwx

// C converts Fahrenheit temperature to Celsius.
func (t TempF) C() TempC {
	return TempC((t - 32.0) / 1.8)
}

// F converts Celsius temperature to Fahrenheit.
func (t TempC) F() TempF {
	return TempF(t*1.8 + 32.0)
}
