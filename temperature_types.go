package libwx

// TempF represents a temperature in degrees Fahrenheit.
type TempF float64

// TempC represents a temperature in degrees Celsius.
type TempC float64

func (t TempF) Unwrap() float64 { return float64(t) }
func (t TempC) Unwrap() float64 { return float64(t) }

type HeatIndexWarning int

const (
	HeatIndexWarningNone = iota
	HeatIndexWarningCaution
	HeatIndexWarningExtremeCaution
	HeatIndexWarningDanger
	HeatIndexWarningExtremeDanger
)
