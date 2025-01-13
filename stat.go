package libwx

import "math"

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

func circularStdDev(x []float64, weights []float64) float64 {
	if weights != nil && len(x) != len(weights) {
		panic("stat: slice length mismatch")
	}

	var aX, aY float64
	if weights != nil {
		var sumW float64
		for i, v := range x {
			w := weights[i]
			sumW += w
			aX += w * math.Cos(v)
			aY += w * math.Sin(v)
		}
		return math.Sqrt(-2 * math.Log(math.Hypot(aY, aX)/sumW))
	} else {
		for _, v := range x {
			aX += math.Cos(v)
			aY += math.Sin(v)
		}
		return math.Sqrt(-2 * math.Log(math.Hypot(aY, aX)/float64(len(x))))
	}
}
