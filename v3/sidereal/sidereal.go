// Copyright 2013 Sonia Keys
// License: MIT

// Sidereal: Chapter 12, Sidereal Time at Greenwich.
package sidereal

import (
	"math"

	"github.com/soniakeys/meeus/v3/base"
	"github.com/soniakeys/meeus/v3/nutation"
	"github.com/soniakeys/unit"
)

// jdToCFrac returns values for use in computing sidereal time at Greenwich.
//
// Cen is centuries from J2000 of the JD at 0h UT of argument jd.  This is
// the value to use for evaluating the IAU sidereal time polynomial.
// DayFrac is the fraction of jd after 0h UT.  It is used to compute the
// final value of sidereal time.
// 计算 T = cen，dayFrac 为此 jd 对应的天的小数
func jdToCFrac(jd float64) (cen, dayFrac float64) {
	j0, f := math.Modf(jd + .5)
	return base.J2000Century(j0 - .5), f
}

// iau82 is a polynomial giving mean sidereal time at Greenwich at 0h UT.
//
// The polynomial is in centuries from J2000.0, as given by JDToCFrac.
// Coefficients are those adopted in 1982 by the International Astronomical
// Union and are given in (12.2) p. 87.
var iau82 = []float64{24110.54841, 8640184.812866, 0.093104, -0.0000062}

// Mean returns mean sidereal time at Greenwich for a given JD.
//
// Computation is by IAU 1982 coefficients.
// The result is in the range [0,86400).
// 计算格林威治 jd 时刻的瞬时平恒星时,化简结果至一天范围之内
func Mean(jd float64) unit.Time {
	return mean(jd).Mod1()
}

// 计算格林威治 jd 时刻的瞬时平恒星时
func mean(jd float64) unit.Time {
	s, f := mean0UT(jd)
	return s + f*1.00273790935
}

// Mean0UT returns mean sidereal time at Greenwich at 0h UT on the given JD.
//
// The result is in the range [0,86400).
// 计算格林威治 0h UT 平恒星时，并化简为单位为秒一天之内的值[0,86400).
func Mean0UT(jd float64) unit.Time {
	s, _ := mean0UT(jd)
	return s.Mod1()
}

// 计算格林威治 0h UT 平恒星时，并返回此 jd 一天中相对于0h的秒数，用于瞬时平恒星时的计算
func mean0UT(jd float64) (sidereal, dayFrac unit.Time) {
	cen, f := jdToCFrac(jd)
	// (12.2) p. 87
	return unit.Time(base.Horner(cen, iau82...)), unit.TimeFromDay(f)
}

// Apparent returns apparent sidereal time at Greenwich for the given JD.
//
// Apparent is mean plus the nutation in right ascension.
//
// The result is in the range [0,86400).
// 计算格林威治瞬时视恒星时
func Apparent(jd float64) unit.Time {
	s := mean(jd)                  // Time
	n := nutation.NutationInRA(jd) // HourAngle
	return (s + n.Time()).Mod1()
}

// Apparent0UT returns apparent sidereal time at Greenwich at 0h UT
// on the given JD.
//
// The result is in the range [0,86400).
// 计算格林威治0h UT视恒星时
func Apparent0UT(jd float64) unit.Time {
	j0, f := math.Modf(jd + .5)
	cen := (j0 - .5 - base.J2000) / 36525
	s := unit.Time(base.Horner(cen, iau82...)) +
		unit.TimeFromDay(f*1.00273790935)
	n := nutation.NutationInRA(j0) // HourAngle
	return (s + n.Time()).Mod1()
}
