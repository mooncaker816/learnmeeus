// Copyright 2013 Sonia Keys
// License: MIT

// Eqtime: Chapter 28, Equation of time.
package eqtime

import (
	"math"

	"github.com/mooncaker816/learnmeeus/v3/base"
	"github.com/mooncaker816/learnmeeus/v3/coord"
	"github.com/mooncaker816/learnmeeus/v3/nutation"
	pp "github.com/mooncaker816/learnmeeus/v3/planetposition"
	"github.com/mooncaker816/learnmeeus/v3/solar"
	"github.com/soniakeys/unit"
)

// E computes the "equation of time" for the given JDE.
// 计算时差
//
// Parameter e must be a planetposition.V87Planet object for Earth obtained
// with planetposition.LoadPlanet.
//
// Result is equation of time as an hour angle.
func E(jde float64, e *pp.V87Planet) unit.HourAngle {
	τ := base.J2000Century(jde) * .1 // J2000儒略日千年数
	L0 := l0(τ)
	// code duplicated from solar.ApparentEquatorialVSOP87 so that
	// we can keep Δψ and cε
	s, β, R := solar.TrueVSOP87(e, jde) // 真太阳黄经
	Δψ, Δε := nutation.Nutation(jde)
	a := unit.AngleFromSec(-20.4898).Div(R) //光行差
	λ := s + Δψ + a                         //视黄经
	ε := nutation.MeanObliquity(jde) + Δε   // 真黄赤交角
	sε, cε := ε.Sincos()
	α, _ := coord.EclToEq(λ, β, sε, cε) // 视赤经
	// (28.1) p. 183
	E := L0 - unit.AngleFromDeg(.0057183) - unit.Angle(α) + Δψ.Mul(cε)
	return unit.HourAngle((E + math.Pi).Mod1() - math.Pi)
}

// (28.2) p. 183
// 太阳平黄经
func l0(τ float64) unit.Angle {
	return unit.AngleFromDeg(base.Horner(τ,
		280.4664567, 360007.6982779, .03032028,
		1./49931, -1./15300, -1./2000000))
}

// ESmart computes the "equation of time" for the given JDE.
// 低精度计算时差
//
// Result is equation of time as an hour angle.
//
// Result is less accurate that E() but the function has the advantage
// of not requiring the V87Planet object.
func ESmart(jde float64) unit.HourAngle {
	ε := nutation.MeanObliquity(jde) // 平黄赤交角
	t := ε.Mul(.5).Tan()
	y := t * t
	T := base.J2000Century(jde)
	L0 := l0(T * .1)
	e := solar.Eccentricity(T) //地球轨道离心率
	M := solar.MeanAnomaly(T)  // 太阳平近点角
	s2L0, c2L0 := L0.Mul(2).Sincos()
	sM := M.Sin()
	// (28.3) p. 185, with double angle identity
	return unit.HourAngle(y*s2L0 - 2*e*sM + 4*e*y*sM*c2L0 -
		y*y*s2L0*c2L0 - 1.25*e*e*M.Mul(2).Sin())
}
