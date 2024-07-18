package libwx

// TempF represents a temperature in degrees Fahrenheit.
type TempF float64

// TempC represents a temperature in degrees Celsius.
type TempC float64

func (t TempF) Unwrap() float64 { return float64(t) }
func (t TempC) Unwrap() float64 { return float64(t) }

type HeatIndexWarning int

const (
	// HeatIndexWarningNone indicates the heat index does not warrant elevated caution.
	HeatIndexWarningNone = iota
	// HeatIndexWarningCaution indicates fatigue is possible with prolonged exposure and activity. Continuing activity could result in heat cramps.
	HeatIndexWarningCaution
	// HeatIndexWarningExtremeCaution indicates heat cramps and heat exhaustion are possible. Continuing activity could result in heat stroke.
	HeatIndexWarningExtremeCaution
	// HeatIndexWarningDangerindicates heat cramps and heat exhaustion are likely; heat stroke is probable with continued activity.
	HeatIndexWarningDanger
	// HeatIndexWarningExtremeDanger indicates heat stroke is imminent.
	HeatIndexWarningExtremeDanger
)
