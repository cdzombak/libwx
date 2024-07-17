package libwx

const (
	ToleranceExact = float64(0.0)
	Tolerance0     = float64(1.0)
	Tolerance1     = float64(0.1)
	Tolerance01    = float64(0.01)
	Tolerance001   = float64(0.001)
	Tolerace1      = Tolerance1 // deprecated; wasa a typo in a previous release
)

func Float64Compare(a, b, tolerance float64) int {
	if a > b+tolerance {
		return 1
	}
	if a < b-tolerance {
		return -1
	}
	return 0
}

func IntCompare(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

func Float64Equal(a, b, tolerance float64) bool {
	return 0 == Float64Compare(a, b, tolerance)
}

func CurriedFloat64Compare(tolerance float64) func(float64, float64) int {
	return func(a, b float64) int {
		return Float64Compare(a, b, tolerance)
	}
}

func CurriedFloat64Equal(tolerance float64) func(float64, float64) bool {
	return func(a, b float64) bool {
		return Float64Equal(a, b, tolerance)
	}
}
