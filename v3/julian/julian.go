// Copyright 2012 Sonia Keys
// License: MIT

// Julian: Chapter 7, Julian day.
//
// Under "General remarks", Meeus describes the INT function as used in the
// book.  In some contexts, math.Floor might be suitable, but I think more
// often, the functions base.FloorDiv or base.FloorDiv64 will be more
// appropriate.  See documentation in package base.
//
// On p. 63, Modified Julian Day is defined.  See constant JMod in package
// base.
//
// See also related functions JulianToGregorian and GregorianToJulian in
// package jm.
package julian

import (
	"math"
	"time"

	"github.com/soniakeys/meeus/v3/base"
)

// CalendarGregorianToJD converts a Gregorian year, month, and day of month
// to Julian day.
//
// Negative years are valid, back to JD 0.  The result is not valid for
// dates before JD 0.
// 格里历日期转儒略日
func CalendarGregorianToJD(y, m int, d float64) float64 {
	switch m {
	case 1, 2:
		y--
		m += 12
	}
	a := base.FloorDiv(y, 100)
	b := 2 - a + base.FloorDiv(a, 4)
	// (7.1) p. 61
	return float64(base.FloorDiv64(36525*(int64(y+4716)), 100)) +
		float64(base.FloorDiv(306*(m+1), 10)+b) + d - 1524.5
}

// CalendarJulianToJD converts a Julian year, month, and day of month to Julian day.
//
// Negative years are valid, back to JD 0.  The result is not valid for
// dates before JD 0.
// 儒略历日期转儒略日
func CalendarJulianToJD(y, m int, d float64) float64 {
	switch m {
	case 1, 2:
		y--
		m += 12
	}
	return float64(base.FloorDiv64(36525*(int64(y+4716)), 100)) +
		float64(base.FloorDiv(306*(m+1), 10)) + d - 1524.5
}

// LeapYearJulian returns true if year y in the Julian calendar is a leap year.
// 儒略历闰年判断
func LeapYearJulian(y int) bool {
	return y%4 == 0
}

// LeapYearGregorian returns true if year y in the Gregorian calendar is a leap year.
// 格里历闰年判断
func LeapYearGregorian(y int) bool {
	return (y%4 == 0 && y%100 != 0) || y%400 == 0
}

// JDToCalendar returns the calendar date for the given jd.
//
// Note that this function returns a date in either the Julian or Gregorian
// Calendar, as appropriate.
// 儒略日转公历日期
// 如果儒略日对应的格里历时间点在1582-10-15 12点之前，则转为儒略历日期
// 如果儒略日对应的格里历时间点在1582-10-15 12点之后，则转为格里历日期
func JDToCalendar(jd float64) (year, month int, day float64) {
	zf, f := math.Modf(jd + .5)
	z := int64(zf)
	a := z
	// if z >= 2299151 { // typo 应该是2299161对应于现实格里历起始日1582-10-15中午12点
	if z >= 2299161 { // typo 应该是2299161对应于现实格里历起始日1582-10-15中午12点
		α := base.FloorDiv64(z*100-186721625, 3652425)
		a = z + 1 + α - base.FloorDiv64(α, 4)
	}
	b := a + 1524
	c := base.FloorDiv64(b*100-12210, 36525)
	d := base.FloorDiv64(36525*c, 100)
	e := int(base.FloorDiv64((b-d)*1e4, 306001))
	// compute return values
	day = float64(int(b-d)-base.FloorDiv(306001*e, 1e4)) + f
	switch e {
	default:
		month = e - 1
	case 14, 15:
		month = e - 13
	}
	switch month {
	default:
		year = int(c) - 4716
	case 1, 2:
		year = int(c) - 4715
	}
	return
}

// jdToCalendarGregorian returns the Gregorian calendar date for the given jd.
//
// Note that it returns a Gregorian date even for dates before the start of
// the Gregorian calendar.  The function is useful when working with Go
// time.Time values because they are always based on the Gregorian calendar.
// 始终转为格里历，忽略儒略历转格里历被扣除的那10天，即把1582-10-15 12点之前的日期也当成格里历算
func jdToCalendarGregorian(jd float64) (year, month int, day float64) {
	zf, f := math.Modf(jd + .5)
	z := int64(zf)
	α := base.FloorDiv64(z*100-186721625, 3652425)
	a := z + 1 + α - base.FloorDiv64(α, 4)
	b := a + 1524
	c := base.FloorDiv64(b*100-12210, 36525)
	d := base.FloorDiv64(36525*c, 100)
	e := int(base.FloorDiv64((b-d)*1e4, 306001))
	// compute return values
	day = float64(int(b-d)-base.FloorDiv(306001*e, 1e4)) + f
	switch e {
	default:
		month = e - 1
	case 14, 15:
		month = e - 13
	}
	switch month {
	default:
		year = int(c) - 4716
	case 1, 2:
		year = int(c) - 4715
	}
	return
}

// JDToTime takes a JD and returns a Go time.Time value.
// 儒略历转为格里历日期对应的 Go Time 类型
func JDToTime(jd float64) time.Time {
	// time.Time is always Gregorian
	y, m, d := jdToCalendarGregorian(jd)
	t := time.Date(y, time.Month(m), 0, 0, 0, 0, 0, time.UTC)
	return t.Add(time.Duration(d * 24 * float64(time.Hour)))
}

// TimeToJD takes a Go time.Time and returns a JD as float64.
//
// Any time zone offset in the time.Time is ignored and the time is
// treated as UTC.
// Go Time 类型转为儒略日（默认日期为格里历）
func TimeToJD(t time.Time) float64 {
	ut := t.UTC()
	y, m, _ := ut.Date()
	d := ut.Sub(time.Date(y, m, 0, 0, 0, 0, 0, time.UTC))
	// time.Time is always Gregorian
	return CalendarGregorianToJD(y, int(m), float64(d)/float64(24*time.Hour))
}

// DayOfWeek determines the day of the week for a given JD.
//
// The value returned is an integer in the range 0 to 6, where 0 represents
// Sunday.  This is the same convention followed in the time package of the
// Go standard library.
// 儒略日0.0为-4712-1-1 12:00 星期一，往后7天一循环就能得出星期几了
func DayOfWeek(jd float64) int {
	return int(jd+1.5) % 7
}

// DayOfYearGregorian computes the day number within the year of the Gregorian
// calendar.
// 格里历年内序数
func DayOfYearGregorian(y, m, d int) int {
	return DayOfYear(y, m, d, LeapYearGregorian(y))
}

// DayOfYearJulian computes the day number within the year of the Julian
// calendar.
// 儒略历年内序数
func DayOfYearJulian(y, m, d int) int {
	return DayOfYear(y, m, d, LeapYearJulian(y))
}

// DayOfYear computes the day number within the year.
//
// This form of the function is not specific to the Julian or Gregorian
// calendar, but you must tell it whether the year is a leap year.
// 输入闰年标识，计算年内序数
func DayOfYear(y, m, d int, leap bool) int {
	k := 2
	if leap {
		k--
	}
	return wholeMonths(m, k) + d
}

// DayOfYearToCalendar returns the calendar month and day for a given
// day of year and leap year status.
// 年内序数求对应的日期
func DayOfYearToCalendar(n int, leap bool) (m, d int) {
	k := 2
	if leap {
		k--
	}
	if n < 32 {
		m = 1
	} else {
		m = (900*(k+n) + 98*275) / 27500
	}
	return m, n - wholeMonths(m, k)
}

// m月之前的所有月份天数之和
func wholeMonths(m, k int) int {
	return 275*m/9 - k*((m+9)/12) - 30
}
