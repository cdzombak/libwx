package libwx

import "testing"

func TestTempF_C(t *testing.T) {
	pairs := []struct {
		input    TempF
		expected TempC
	}{
		{TempF(32), TempC(0)},
		{TempF(212), TempC(100)},
		{TempF(98.6), TempC(37)},
		{TempF(-40), TempC(-40)},
	}

	for _, pair := range pairs {
		if Float64Compare(float64(pair.expected), float64(pair.input.C()), Tolerance001) != 0 {
			t.Errorf("for input %v: expected %v, got %v", pair.input, pair.expected, pair.input.C())
		}
	}
}

func TestTempC_F(t *testing.T) {
	pairs := []struct {
		input    TempC
		expected TempF
	}{
		{TempC(0), TempF(32)},
		{TempC(100), TempF(212)},
		{TempC(37), TempF(98.6)},
		{TempC(-40), TempF(-40)},
	}

	for _, pair := range pairs {
		if Float64Compare(float64(pair.expected), float64(pair.input.F()), Tolerance001) != 0 {
			t.Errorf("for input %v: expected %v, got %v", pair.input, pair.expected, pair.input.F())
		}
	}
}
