package main

import (
	"math"
)

func mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	var sum float64
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func diff(a []float64, b float64) []float64 {
	var output = make([]float64, len(a))
	for i, x := range a {
		output[i] = x - b
	}
	return output
}

func dotProd(x, y []float64) float64 {
	var sum float64
	for i := range x {
		sum += (x[i] * y[i])
	}
	return sum
}

func coorelation(x, y []float64) float64 {
	xMean := mean(x)
	yMean := mean(y)

	xMeanDiff := diff(x, xMean)
	yMeanDiff := diff(y, yMean)

	ab := dotProd(xMeanDiff, yMeanDiff)
	a2 := dotProd(xMeanDiff, xMeanDiff)
	b2 := dotProd(yMeanDiff, yMeanDiff)
	if a2 == 0 || b2 == 0 {
		return 0
	}

	return ab / math.Sqrt(a2*b2)
}
