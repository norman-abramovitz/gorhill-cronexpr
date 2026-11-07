/*!
 * Copyright 2026 Norman Abramovitz
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_prev.go
 * License: pick the one which suits you :
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

package cronexpr

/******************************************************************************/

import (
	"sort"
	"time"
)

// The mirror of cronexpr_next.go: each step finds the last listed value
// before the current one, or moves back to the previous larger unit.

/******************************************************************************/

// prevWall returns the last wall-clock time before `fromTime` which matches
// the cron expression `expr`, with wall-clock times carried in UTC as for
// nextWall.
func (expr *Expression) prevWall(fromTime time.Time) time.Time {
	// As in nextWall: move back as soon as a field does not match.

	// year
	v := fromTime.Year()
	i := sort.SearchInts(expr.yearList, v)
	if i == len(expr.yearList) || v != expr.yearList[i] {
		return expr.prevYear(fromTime)
	}
	// month
	v = int(fromTime.Month())
	i = sort.SearchInts(expr.monthList, v)
	if i == len(expr.monthList) || v != expr.monthList[i] {
		return expr.prevMonth(fromTime)
	}

	days := expr.calculateActualDaysOfMonth(fromTime.Year(), int(fromTime.Month()))

	// day of month
	v = fromTime.Day()
	i = sort.SearchInts(days, v)
	if i == len(days) || v != days[i] {
		return expr.prevDayOfMonth(fromTime)
	}
	// hour
	v = fromTime.Hour()
	i = sort.SearchInts(expr.hourList, v)
	if i == len(expr.hourList) || v != expr.hourList[i] {
		return expr.prevHour(fromTime)
	}
	// minute
	v = fromTime.Minute()
	i = sort.SearchInts(expr.minuteList, v)
	if i == len(expr.minuteList) || v != expr.minuteList[i] {
		return expr.prevMinute(fromTime)
	}

	// The second matches or not: the previous second is the answer
	return expr.prevSecond(fromTime)
}

/******************************************************************************/

// last returns the last value of a sorted list.
func last(list []int) int {
	return list[len(list)-1]
}

/******************************************************************************/

func (expr *Expression) prevYear(t time.Time) time.Time {
	// Find index of the last item in list smaller than candidate year
	i := sort.SearchInts(expr.yearList, t.Year()) - 1
	if i < 0 {
		return time.Time{}
	}
	// Year changed, need to recalculate actual days of month
	month := last(expr.monthList)
	days := expr.calculateActualDaysOfMonth(expr.yearList[i], month)
	if len(days) == 0 {
		return expr.prevMonth(time.Date(expr.yearList[i], time.Month(month), 1, 0, 0, 0, 0, t.Location()))
	}
	return time.Date(
		expr.yearList[i],
		time.Month(month),
		last(days),
		last(expr.hourList),
		last(expr.minuteList),
		last(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevMonth(t time.Time) time.Time {
	// Find index of the last item in list smaller than candidate month
	i := sort.SearchInts(expr.monthList, int(t.Month())) - 1
	if i < 0 {
		return expr.prevYear(t)
	}
	// Month changed, need to recalculate actual days of month
	days := expr.calculateActualDaysOfMonth(t.Year(), expr.monthList[i])
	if len(days) == 0 {
		return expr.prevMonth(time.Date(t.Year(), time.Month(expr.monthList[i]), 1, 0, 0, 0, 0, t.Location()))
	}
	return time.Date(
		t.Year(),
		time.Month(expr.monthList[i]),
		last(days),
		last(expr.hourList),
		last(expr.minuteList),
		last(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevDayOfMonth(t time.Time) time.Time {
	// The days of t's month: the caller has matched its year and month
	days := expr.calculateActualDaysOfMonth(t.Year(), int(t.Month()))
	// Find index of the last item in list smaller than candidate day of month
	i := sort.SearchInts(days, t.Day()) - 1
	if i < 0 {
		return expr.prevMonth(t)
	}
	return time.Date(
		t.Year(),
		t.Month(),
		days[i],
		last(expr.hourList),
		last(expr.minuteList),
		last(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevHour(t time.Time) time.Time {
	// Find index of the last item in list smaller than candidate hour
	i := sort.SearchInts(expr.hourList, t.Hour()) - 1
	if i < 0 {
		return expr.prevDayOfMonth(t)
	}
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		expr.hourList[i],
		last(expr.minuteList),
		last(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevMinute(t time.Time) time.Time {
	// Find index of the last item in list smaller than candidate minute
	i := sort.SearchInts(expr.minuteList, t.Minute()) - 1
	if i < 0 {
		return expr.prevHour(t)
	}
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		expr.minuteList[i],
		last(expr.secondList),
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) prevSecond(t time.Time) time.Time {
	// prevSecond() assumes all other fields are exactly matched
	// to the cron expression

	// Find index of the last item in list smaller than candidate second
	i := sort.SearchInts(expr.secondList, t.Second()) - 1
	if i < 0 {
		return expr.prevMinute(t)
	}
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		expr.secondList[i],
		0,
		t.Location())
}
