package libwx

import "math"

// borrowed from gonum:
func circularMean(x, weights []float64) float64 {
	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}

	var aX, aY float64
	if weights != nil {
		for i, v := range x {
			aX += weights[i] * math.Cos(v)
			aY += weights[i] * math.Sin(v)
		}
	} else {
		for _, v := range x {
			aX += math.Cos(v)
			aY += math.Sin(v)
		}
	}

	return math.Atan2(aY, aX)
}

// from my (open) gonum PR which doesn't (yet) support weights:
func circularStdDev(x []float64) float64 {
	var aX, aY float64
	for _, v := range x {
		aX += math.Cos(v)
		aY += math.Sin(v)
	}
	return math.Sqrt(-2 * math.Log(math.Sqrt(aY*aY+aX*aX)/float64(len(x))))
}
