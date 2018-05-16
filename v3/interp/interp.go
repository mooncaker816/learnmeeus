// Copyright 2012 Sonia Keys
// License: MIT

// Interp: Chapter 3, Interpolation.
//
// Len3 and Len5 types
//
// These types allow interpolation from a table of equidistant x values
// and corresponding y values.  Since the x values are equidistant,
// only the first and last values are supplied as arguments to the
// constructors.  The interior x values are implicit.  All y values must be
// supplied however.  They are passed as a slice, and the length of y is fixed.
// For Len3 it must be 3 and for Len5 it must be 5.
//
// For these Len3 and Len5 functions, Meeus notes the importance of choosing
// the 3 or 5 rows of a larger table that will minimize the interpolating
// factor n.  He does not provide algorithms for doing this however.
//
// For an example of a selection function, see Len3ForInterpolateX.  This
// was useful for computing Delta T.
package interp

import (
	"errors"
	"math"

	"github.com/soniakeys/meeus/v3/base"
)

// Error values returned by functions and methods in this package.
// Defined here to help testing for specific errors.
var (
	ErrorNot3            = errors.New("Argument y must be length 3")
	ErrorNot4            = errors.New("Argument y must be length 4")
	ErrorNot5            = errors.New("Argument y must be length 5")
	ErrorNoXRange        = errors.New("Argument x3 (or x5) cannot equal x1")
	ErrorNOutOfRange     = errors.New("Interpolating factor n must be in range -1 to 1")
	ErrorXOutOfRange     = errors.New("Argument x outside of range x1 to x3 (or x5)")
	ErrorNoExtremum      = errors.New("No extremum in table")
	ErrorExtremumOutside = errors.New("Extremum falls outside of table")
	ErrorZeroOutside     = errors.New("Zero falls outside of table")
	ErrorNoConverge      = errors.New("Failure to converge")
)

// Len3 allows second difference interpolation.
// 等距三点插值结构
type Len3 struct {
	x1, x3             float64   //x1,x3分别为起始点和终止点，无需给出x2，因为等距
	y                  []float64 //y为x1,x2,x3对应的值的序列
	a, b, c            float64   //a=y2-y1, b=y3-y2, c=b-a=y3+y1-2y2
	abSum, xSum, xDiff float64   //计数插值的中间变量
}

// NewLen3 prepares a Len3 object from a table of three rows of x and y values.
//
// X values must be equally spaced, so only the first and last are supplied.
// X1 must not equal x3.  Y must be a slice of three y values.
// 根据上述定义，创建三点插值结构
func NewLen3(x1, x3 float64, y []float64) (*Len3, error) {
	if len(y) != 3 {
		return nil, ErrorNot3
	}
	if x3 == x1 {
		return nil, ErrorNoXRange
	}
	d := &Len3{
		x1: x1,
		x3: x3,
		y:  append([]float64{}, y...),
	}
	// differences. (3.1) p. 23
	d.a = y[1] - y[0]
	d.b = y[2] - y[1]
	d.c = d.b - d.a
	// other intermediate values
	d.abSum = d.a + d.b
	d.xSum = x3 + x1
	d.xDiff = x3 - x1
	return d, nil
}

// Len3ForInterpolateX is a special purpose Len3 constructor.
//
// Like NewLen3, it takes a table of x and y values, but it is not limited
// to tables of 3 rows.  An X value is also passed that represents the
// interpolation target x value.  Len3ForInterpolateX will locate the
// appropriate three rows of the table for interpolating for x, and initialize
// the Len3 object for those rows.
//
//	x is the target for interpolation
//	x1 is the x value corresponding to the first y value of the table.
//	xn is the x value corresponding to the last y value of the table.
//	y is all y values in the table.  len(y) should be >= 3.
//
// 给定n个点，但是我们只需选取离目标点x最接近的三个点来做三点插值，
// 此时可用以下函数来自动选择最优三点，来构造三点插值
// 同样，前提是n个点等距，且与y一一对应
func Len3ForInterpolateX(x, x1, xn float64, y []float64) (*Len3, error) {
	if len(y) > 3 {
		interval := (xn - x1) / float64(len(y)-1)
		if interval == 0 {
			return nil, ErrorNoXRange
		}
		nearestX := int((x-x1)/interval + .5)
		if nearestX < 1 {
			nearestX = 1
		} else if nearestX > len(y)-2 {
			nearestX = len(y) - 2
		}
		y = y[nearestX-1 : nearestX+2]
		xn = x1 + float64(nearestX+1)*interval
		x1 = x1 + float64(nearestX-1)*interval
	}
	return NewLen3(x1, xn, y)
}

// InterpolateX interpolates for a given x value.
// 计算插值因子n，调用非严格插值计算
func (d *Len3) InterpolateX(x float64) (y float64) {
	n := (2*x - d.xSum) / d.xDiff
	return d.InterpolateN(n)
}

// InterpolateXStrict interpolates for a given x value,
// restricting x to the range x1 to x3 given to the constructor NewLen3.
// 计算插值因子n，调用严格插值计算
func (d *Len3) InterpolateXStrict(x float64) (y float64, err error) {
	n := (2*x - d.xSum) / d.xDiff
	y, err = d.InterpolateNStrict(n)
	if err == ErrorNOutOfRange {
		err = ErrorXOutOfRange
	}
	return
}

// InterpolateN interpolates for a given interpolating factor n.
//
// This is interpolation formula (3.3)
//
// The interpolation factor n is x-x2 in units of the tabular x interval.
// (See Meeus p. 24.)
// 非严格插值计算，不用保证目标点插值因子绝对值小于等于1，
// 即不用保证离我们所选三点中点距离小于一个步长
func (d *Len3) InterpolateN(n float64) (y float64) {
	return d.y[1] + n*.5*(d.abSum+n*d.c)
}

// InterpolateNStrict interpolates for a given interpolating factor n.
//
// N is restricted to the range [-1..1] corresponding to the range x1 to x3
// given to the constructor NewLen3.
// 严格插值计算，必须保证目标点插值因子绝对值小于等于1，
// 即必须保证离我们所选三点中点距离小于一个步长
func (d *Len3) InterpolateNStrict(n float64) (y float64, err error) {
	if n < -1 || n > 1 {
		return 0, ErrorNOutOfRange
	}
	return d.InterpolateN(n), nil
}

// Extremum returns the x and y values at the extremum.
//
// Results are restricted to the range of the table given to the constructor
// NewLen3.
func (d *Len3) Extremum() (x, y float64, err error) {
	if d.c == 0 {
		return 0, 0, ErrorNoExtremum
	}
	n := d.abSum / (-2 * d.c) // (3.5), p. 25
	if n < -1 || n > 1 {
		return 0, 0, ErrorExtremumOutside
	}
	x = .5 * (d.xSum + d.xDiff*n)          // 根据实际步长得出极值点x
	y = d.y[1] - (d.abSum*d.abSum)/(8*d.c) // (3.4), p. 25
	return x, y, nil
}

// Len3Zero finds a zero of the quadratic function represented by the table.
//
// That is, it returns an x value that yields y=0.
//
// Argument strong switches between two strategies for the estimation step.
// when iterating to converge on the zero.
//
// Strong=false specifies a quick and dirty estimate that works well
// for gentle curves, but can work poorly or fail on more dramatic curves.
//
// Strong=true specifies a more sophisticated and thus somewhat more
// expensive estimate.  However, if the curve has quick changes, This estimate
// will converge more reliably and in fewer steps, making it a better choice.
//
// Results are restricted to the range of the table given to the constructor
// NewLen3.
// strong 为考虑修正量的迭代方式，更为精确
func (d *Len3) Zero(strong bool) (x float64, err error) {
	var f iterFunc
	if strong {
		// (3.7), p. 27
		f = func(n0 float64) float64 {
			return n0 - (2*d.y[1]+n0*(d.abSum+d.c*n0))/(d.abSum+2*d.c*n0)
		}
	} else {
		// (3.6), p. 26
		f = func(n0 float64) float64 {
			return -2 * d.y[1] / (d.abSum + d.c*n0)
		}
	}
	n0, ok := iterate(0, f)
	if !ok {
		return 0, ErrorNoConverge
	}
	if n0 > 1 || n0 < -1 {
		return 0, ErrorZeroOutside
	}
	return .5 * (d.xSum + d.xDiff*n0), nil // success
}

type iterFunc func(n0 float64) (n1 float64)

func iterate(n0 float64, f iterFunc) (n1 float64, ok bool) {
	for limit := 0; limit < 50; limit++ {
		n1 = f(n0)
		if math.IsInf(n1, 0) || math.IsNaN(n1) {
			break // failure to converge
		}
		if math.Abs((n1-n0)/n0) < 1e-15 {
			return n1, true // success
		}
		n0 = n1
	}
	return 0, false // failure to converge
}

// Len4Half interpolates a center value from a table of four rows.
func Len4Half(y []float64) (float64, error) {
	if len(y) != 4 {
		return 0, ErrorNot4
	}
	// (3.12) p. 32
	return (9*(y[1]+y[2]) - y[0] - y[3]) / 16, nil
}

// Len5 allows fourth difference interpolation.
// 五点等距插值结构
type Len5 struct {
	x1, x5      float64   // x1为起始点，x5为终止点
	y           []float64 // y 为x1,x2,x3,x4,x5的一一映射
	a, b, c, d  float64   // a=y2-y1, b=y3-y2, c=y4-y3, d=y5-y4
	e, f, g     float64   // e=b-a, f=c-b, g=d-c
	h, j, k     float64   // h=f-e, j=g-f, k=j-h
	y3          float64   // y3为中间点的 y 值
	xSum, xDiff float64   // xSum=x1+x5, xDiff=x5-x1
	interpCoeff []float64 // 插值函数对应插值因子 n 的各项系数（0-4）
}

// NewLen5 prepares a Len5 object from a table of five rows of x and y values.
//
// X values must be equally spaced, so only the first and last are supplied.
// X1 must not equal x5.  Y must be a slice of five y values.
// 构造5点插值结构
func NewLen5(x1, x5 float64, y []float64) (*Len5, error) {
	if len(y) != 5 {
		return nil, ErrorNot5
	}
	if x5 == x1 {
		return nil, ErrorNoXRange
	}
	d := &Len5{
		x1: x1,
		x5: x5,
		y:  append([]float64{}, y...),
		y3: y[2],
	}
	// differences
	d.a = y[1] - y[0]
	d.b = y[2] - y[1]
	d.c = y[3] - y[2]
	d.d = y[4] - y[3]

	d.e = d.b - d.a
	d.f = d.c - d.b
	d.g = d.d - d.c

	d.h = d.f - d.e
	d.j = d.g - d.f

	d.k = d.j - d.h
	// other intermediate values
	d.xSum = x5 + x1
	d.xDiff = x5 - x1
	d.interpCoeff = []float64{ // (3.8) p. 28
		d.y3,
		(d.b+d.c)/2 - (d.h+d.j)/12,
		d.f/2 - d.k/24,
		(d.h + d.j) / 12,
		d.k / 24,
	}
	return d, nil
}

// InterpolateX interpolates for a given x value.
func (d *Len5) InterpolateX(x float64) (y float64) {
	n := (4*x - 2*d.xSum) / d.xDiff
	return d.InterpolateN(n)
}

// InterpolateXStrict interpolates for a given x value,
// restricting x to the range x1 to x5 given to the the constructor NewLen5.
func (d *Len5) InterpolateXStrict(x float64) (y float64, err error) {
	n := (4*x - 2*d.xSum) / d.xDiff
	y, err = d.InterpolateNStrict(n)
	if err == ErrorNOutOfRange {
		err = ErrorXOutOfRange
	}
	return
}

// InterpolateN interpolates for a given interpolating factor n.
//
// The interpolation factor n is x-x3 in units of the tabular x interval.
// (See Meeus p. 28.)
// Horner 为工具函数，求解多项式之和，interpCoeff为多项式系数
func (d *Len5) InterpolateN(n float64) (y float64) {
	return base.Horner(n, d.interpCoeff...)
}

// InterpolateNStrict interpolates for a given interpolating factor n.
//
// N is restricted to the range [-1..1].  This is only half the range given
// to the constructor NewLen5, but is the recommendation given on p. 31.
func (d *Len5) InterpolateNStrict(n float64) (y float64, err error) {
	if n < -1 || n > 1 {
		return 0, ErrorNOutOfRange
	}
	return base.Horner(n, d.interpCoeff...), nil
}

// Extremum returns the x and y values at the extremum.
//
// Results are restricted to the range of the table given to the constructor
// NewLen5.  (Meeus actually recommends restricting the range to one unit of
// the tabular interval, but that seems a little harsh.)
// 5点插值极值
func (d *Len5) Extremum() (x, y float64, err error) {
	// (3.9) p. 29
	nCoeff := []float64{
		6*(d.b+d.c) - d.h - d.j,
		0,
		// 3 * (d.h + d.k), // 不应该是 d.h+d.j 吗？
		3 * (d.h + d.j), // 不应该是 d.h+d.j 吗？
		2 * d.k,
	}
	den := d.k - 12*d.f
	if den == 0 {
		return 0, 0, ErrorExtremumOutside
	}
	n0, ok := iterate(0, func(n0 float64) float64 {
		return base.Horner(n0, nCoeff...) / den
	})
	if !ok {
		return 0, 0, ErrorNoConverge
	}
	if n0 < -2 || n0 > 2 {
		return 0, 0, ErrorExtremumOutside
	}
	x = .5*d.xSum + .25*d.xDiff*n0
	y = base.Horner(n0, d.interpCoeff...)
	return x, y, nil
}

// Len5Zero finds a zero of the quartic function represented by the table.
//
// That is, it returns an x value that yields y=0.
//
// Argument strong switches between two strategies for the estimation step.
// when iterating to converge on the zero.
//
// Strong=false specifies a quick and dirty estimate that works well
// for gentle curves, but can work poorly or fail on more dramatic curves.
//
// Strong=true specifies a more sophisticated and thus somewhat more
// expensive estimate.  However, if the curve has quick changes, This estimate
// will converge more reliably and in fewer steps, making it a better choice.
//
// Results are restricted to the range of the table given to the constructor
// NewLen5.
// strong 为带修正模式
func (d *Len5) Zero(strong bool) (x float64, err error) {
	var f iterFunc
	if strong {
		// (3.11), p. 29
		M := d.k / 24
		N := (d.h + d.j) / 12
		P := d.f/2 - M
		Q := (d.b+d.c)/2 - N
		numCoeff := []float64{d.y3, Q, P, N, M}
		denCoeff := []float64{Q, 2 * P, 3 * N, 4 * M}
		f = func(n0 float64) float64 {
			return n0 -
				base.Horner(n0, numCoeff...)/base.Horner(n0, denCoeff...)
		}
	} else {
		// (3.10), p. 29
		numCoeff := []float64{
			-24 * d.y3,
			0,
			d.k - 12*d.f,
			-2 * (d.h + d.j),
			-d.k,
		}
		den := 12*(d.b+d.c) - 2*(d.h+d.j)
		f = func(n0 float64) float64 {
			return base.Horner(n0, numCoeff...) / den
		}
	}
	n0, ok := iterate(0, f)
	if !ok {
		return 0, ErrorNoConverge
	}
	if n0 > 2 || n0 < -2 {
		return 0, ErrorZeroOutside
	}
	x = .5*d.xSum + .25*d.xDiff*n0
	return x, nil
}

// Lagrange performs interpolation with unequally-spaced abscissae.
//
// Given a table of X and Y values, interpolate a new y value for argument x.
//
// X values in the table do not have to be equally spaced; they do not even
// have to be in order.  They must however, be distinct.
// table 中包含了 n 个点且xi 必须互异,x 为目标插值点
func Lagrange(x float64, table []struct{ X, Y float64 }) (y float64) {
	// method of BASIC program, p. 33.
	sum := 0.
	for i := range table {
		xi := table[i].X
		prod := 1.
		for j := range table {
			if i != j {
				xj := table[j].X
				prod *= (x - xj) / (xi - xj)
			}
		}
		sum += table[i].Y * prod
	}
	return sum
}

// LagrangePoly uses the formula of Lagrange to produce an interpolating
// polynomial.
//
// X values in the table do not have to be equally spaced; they do not even
// have to be in order.  They must however, be distinct.
//
// The returned polynomial will be of degree n-1 where n is the number of rows
// in the table.  It can be evaluated for x using common.Horner.
func LagrangePoly(table []struct{ X, Y float64 }) []float64 {
	// Method not fully described by Meeus, but needed for numerical solution
	// to Example 3.g.
	sum := make([]float64, len(table))
	prod := make([]float64, len(table))
	last := len(table) - 1
	for i := range table {
		xi := table[i].X
		yi := table[i].Y
		prod[last] = 1
		den := 1.
		n := last
		for j := range table {
			if i != j {
				xj := table[j].X
				prod[n-1] = prod[n] * -xj
				for k := n; k < last; k++ {
					prod[k] -= prod[k+1] * xj
				}
				n--
				den *= (xi - xj)
			}
		}
		for j, pj := range prod {
			sum[j] += yi * pj / den
		}
	}
	return sum
}
