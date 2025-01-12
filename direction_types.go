package libwx

// Degree represents direction in (angular) degrees.
type Degree float64

func (d Degree) Unwrap() float64 { return float64(d) }
