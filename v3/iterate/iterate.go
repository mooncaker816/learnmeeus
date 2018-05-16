// Copyright 2013 Sonia Keys
// License: MIT

// Iterate: Chapter 5, Iteration.
//
// This package is best considered illustrative.  While the functions are
// usable, they are minimal in showing the points of the chapter text.  More
// robust functions would handle more cases of overflow, loss of precision,
// and divergence.
package iterate

import (
	"errors"
	"math"
)

// BetterFunc is a convience type definition.
type BetterFunc func(float64) float64

// DecimalPlaces iterates to a fixed number of decimal places.
//
// Inputs are an improvement function, a starting value, the number of
// decimal places desired in the result, and an iteration limit.
// better 为迭代公式，start 为起始点，places 为视作迭代结束的最大精度值的小数点位数
// maxIterations 为迭代最大次数
func DecimalPlaces(better BetterFunc, start float64, places, maxIterations int) (float64, error) {
	d := math.Pow(10, float64(-places))
	for i := 0; i < maxIterations; i++ {
		n := better(start)
		if math.Abs(n-start) < d {
			return n, nil
		}
		start = n
	}
	return 0, errors.New("Maximum iterations reached")
}

// FullPrecison iterates to (nearly) the full precision of a float64.
//
// To allow for a little bit of floating point jitter, FullPrecision iterates
// to 15 significant figures, which is the maximum number of full significant
// figures representable in a float64, but still a couple of bits shy of the
// full representable precision.
// 和 DecimalPlaces 功能类似，只不过默认迭代结束的标志为小于float64的最大精度15
func FullPrecision(better BetterFunc, start float64, maxIterations int) (float64, error) {
	for i := 0; i < maxIterations; i++ {
		n := better(start)
		if math.Abs((n-start)/n) < 1e-15 {
			return n, nil
		}
		start = n
	}
	return 0, errors.New("Maximum iterations reached")
}

// RootFunc is a convience type definition.
type RootFunc func(float64) float64

// BinaryRoot finds a root between given bounds by binary search.
//
// Inputs are a function on x and the bounds on x.  A root must exist between
// the given bounds, otherwise the result is not meaningful.
// 因为float64 的小数bit位最多为52位(10进制为15位有效数字)，
// 即每次迭代理论上能提升一位精度，所以最多52次就应该跳出迭代
func BinaryRoot(f RootFunc, lower, upper float64) float64 {
	yLower := f(lower)
	var mid float64
	for j := 0; j < 52; j++ {
		mid = (lower + upper) / 2
		yMid := f(mid)
		if yMid == 0 {
			break
		}
		if math.Signbit(yLower) == math.Signbit(yMid) { // 与低值同号，替代低值
			lower = mid
			yLower = yMid
		} else {
			upper = mid //与高值同号，替换高值
		}
	}
	return mid
}
