// Copyright 2013 Sonia Keys
// License: MIT

// Stellar: Chapter 56, Stellar Magnitudes.
package stellar

import (
	"math"

	"github.com/soniakeys/unit"
)

// Sum returns the combined apparent magnitude of two stars.
// 2个恒星的合星等
func Sum(m1, m2 float64) float64 {
	x := .4 * (m2 - m1)
	return m2 - 2.5*math.Log10(math.Pow(10, x)+1)
}

// SumN returns the combined apparent magnitude of a number of stars.
// n个恒星的合星等
func SumN(m ...float64) float64 {
	var s float64
	for _, mi := range m {
		s += math.Pow(10, -.4*mi)
	}
	return -2.5 * math.Log10(s)
}

// Ratio returns the brightness ratio of two stars.
// 天体1和天体2的视亮度比
//
// Arguments m1, m2 are apparent magnitudes.
func Ratio(m1, m2 float64) float64 {
	x := .4 * (m2 - m1)
	return math.Pow(10, x)
}

// Difference returns the difference in apparent magnitude of two stars
// given their brightness ratio.
// 通过亮度比计算两个恒星的星等差
func Difference(ratio float64) float64 {
	return 2.5 * math.Log10(ratio)
}

// AbsoluteByParallax returns absolute magnitude given annual parallax.
// 秒差距+视星等计算绝对星等
// 恒星的绝对星等是指：我们位于距恒星10(秒差距)的地方得到的恒星视星等
//
// Argument m is apparent magnitude, π is annual parallax.
func AbsoluteByParallax(m float64, π unit.Angle) float64 {
	return m + 5 + 5*math.Log10(π.Sec())
}

// AbsoluteByDistance returns absolute magnitude given distance.
// 距离+视星等计算绝对星等
//
// Argument m is apparent magnitude, d is distance in parsecs.
func AbsoluteByDistance(m, d float64) float64 {
	return m + 5 - 5*math.Log10(d)
}
