// Copyright 2013 Sonia Keys
// License: MIT

// Pluto: Chapter 37, Pluto.
package pluto

import (
	"math"

	"github.com/mooncaker816/learnmeeus/v3/base"
	"github.com/mooncaker816/learnmeeus/v3/elliptic"
	pp "github.com/mooncaker816/learnmeeus/v3/planetposition"
	"github.com/soniakeys/unit"
)

// Heliocentric returns J2000 heliocentric coordinates of Pluto.
//  J2000冥王星日心黄道坐标
//
// Results l, b are solar longitude and latitude in radians.
// Result r is distance in AU.
func Heliocentric(jde float64) (l, b unit.Angle, r float64) {
	T := base.J2000Century(jde)
	J := unit.AngleFromDeg(34.35 + 3034.9057*T)
	S := unit.AngleFromDeg(50.08 + 1222.1138*T)
	P := unit.AngleFromDeg(238.96 + 144.96*T)
	for i := range t37 {
		t := &t37[i]
		sα, cα := (J.Mul(t.i) + S.Mul(t.j) + P.Mul(t.k)).Sincos()
		l += t.lA.Mul(sα) + t.lB.Mul(cα)
		b += t.bA.Mul(sα) + t.bB.Mul(cα)
		r += t.rA*sα + t.rB*cα
	}
	l += unit.AngleFromDeg(238.958116 + 144.96*T)
	b -= unit.AngleFromDeg(3.908239)
	r += 40.7241346
	return
}

// Astrometric returns J2000 astrometric coordinates of Pluto.
//  J2000冥王星地心赤道坐标
func Astrometric(jde float64, e *pp.V87Planet) (α unit.RA, δ unit.Angle) {
	const sε, cε = base.SOblJ2000, base.COblJ2000
	f := func(jde float64) (x, y, z float64) {
		l, b, r := Heliocentric(jde)
		sl, cl := l.Sincos()
		sb, cb := b.Sincos()
		// (37.1) p. 264
		x = r * cl * cb
		y = r * (sl*cb*cε - sb*sε)
		z = r * (sl*cb*sε + sb*cε)
		return
	}
	α, δ, _ = elliptic.AstrometricJ2000(f, jde, e)
	return
}

func init() {
	for i := range t37 {
		t := &t37[i]
		t.lA *= math.Pi / 180
		t.lB *= math.Pi / 180
		t.bA *= math.Pi / 180
		t.bB *= math.Pi / 180
	}
}

var t37 = []struct {
	i, j, k float64
	lA, lB  unit.Angle
	bA, bB  unit.Angle
	rA, rB  float64
}{
	{0, 0, 1, -19.799805, 19.850055, -5.452852, -14.974862, 6.6865439, 6.8951812},
	{0, 0, 2, .897144, -4.954829, 3.527812, 1.67279, -1.1827535, -.0332538},
	{0, 0, 3, .611149, 1.211027, -1.050748, .327647, .1593179, -.143889},
	{0, 0, 4, -.341243, -.189585, .17869, -.292153, -.0018444, .048322},
	{0, 0, 5, .129287, -.034992, .01865, .10034, -.0065977, -.0085431},
	{0, 0, 6, -.038164, .030893, -.030697, -.025823, .0031174, -.0006032},
	{0, 1, -1, .020442, -.009987, .004878, .011248, -.0005794, .0022161},
	{0, 1, 0, -.004063, -.005071, .000226, -.000064, .0004601, .0004032},
	{0, 1, 1, -.006016, -.003336, .00203, -.000836, -.0001729, .0000234},
	{0, 1, 2, -.003956, .003039, .000069, -.000604, -.0000415, .0000702},
	{0, 1, 3, -.000667, .003572, -.000247, -.000567, .0000239, .0000723},
	{0, 2, -2, .001276, .000501, -.000057, .000001, .0000067, -.0000067},
	{0, 2, -1, .001152, -.000917, -.000122, .000175, .0001034, -.0000451},
	{0, 2, 0, .00063, -.001277, -.000049, -.000164, -.0000129, .0000504},
	{1, -1, 0, .002571, -.000459, -.000197, .000199, .000048, -.0000231},
	{1, -1, 1, .000899, -.001449, -.000025, .000217, .0000002, -.0000441},
	{1, 0, -3, -.001016, .001043, .000589, -.000248, -.0003359, .0000265},
	{1, 0, -2, -.002343, -.001012, -.000269, .000711, .0007856, -.0007832},
	{1, 0, -1, .007042, .000788, .000185, .000193, .0000036, .0045763},
	{1, 0, 0, .001199, -.000338, .000315, .000807, .0008663, .0008547},
	{1, 0, 1, .000418, -.000067, -.00013, -.000043, -.0000809, -.0000769},
	{1, 0, 2, .00012, -.000274, .000005, .000003, .0000263, -.0000144},
	{1, 0, 3, -.00006, -.000159, .000002, .000017, -.0000126, .0000032},
	{1, 0, 4, -.000082, -.000029, .000002, .000005, -.0000035, -.0000016},
	{1, 1, -3, -.000036, -.000029, .000002, .000003, -.0000019, -.0000004},
	{1, 1, -2, -.00004, .000007, .000003, .000001, -.0000015, .0000008},
	{1, 1, -1, -.000014, .000022, .000002, -.000001, -.0000004, .0000012},
	{1, 1, 0, .000004, .000013, .000001, -.000001, .0000005, .0000006},
	{1, 1, 1, .000005, .000002, 0, -.000001, .0000003, .0000001},
	{1, 1, 3, -.000001, 0, 0, 0, .0000006, -.0000002},
	{2, 0, -6, .000002, 0, 0, -.000002, .0000002, .0000002},
	{2, 0, -5, -.000004, .000005, .000002, .000002, -.0000002, -.0000002},
	{2, 0, -4, .000004, -.000007, -.000007, 0, .0000014, .0000013},
	{2, 0, -3, .000014, .000024, .00001, -.000008, -.0000063, .0000013},
	{2, 0, -2, -.000049, -.000034, -.000003, .00002, .0000136, -.0000236},
	{2, 0, -1, .000163, -.000048, .000006, .000005, .0000273, .0001065},
	{2, 0, 0, .000009, -.000024, .000014, .000017, .0000251, .0000149},
	{2, 0, 1, -.000004, .000001, -.000002, 0, -.0000025, -.0000009},
	{2, 0, 2, -.000003, .000001, 0, 0, .0000009, -.0000002},
	{2, 0, 3, .000001, .000003, 0, 0, -.0000008, .0000007},
	{3, 0, -2, -.000003, -.000001, 0, .000001, .0000002, -.000001},
	{3, 0, -1, .000005, -.000003, 0, 0, .0000019, .0000035},
	{3, 0, 0, 0, 0, .000001, 0, .000001, .0000003},
}
