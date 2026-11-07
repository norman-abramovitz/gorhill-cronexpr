/*!
 * Copyright 2013 Raymond Hill
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_next.go
 * Version: 1.0
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

/******************************************************************************/

var dowNormalizedOffsets = [][]int{
	{1, 8, 15, 22, 29},
	{2, 9, 16, 23, 30},
	{3, 10, 17, 24, 31},
	{4, 11, 18, 25},
	{5, 12, 19, 26},
	{6, 13, 20, 27},
	{7, 14, 21, 28},
}

/******************************************************************************/

// nextWall returns the first wall-clock time after `fromTime` which matches
// the cron expression `expr`. Wall-clock times are carried in UTC, where
// every time of every day exists once, so the steps below never meet a
// daylight saving change; nextExisting maps them to the zone.
func (expr *Expression) nextWall(fromTime time.Time) time.Time {
	// Since expr.nextSecond()-expr.nextMonth() expects that the
	// supplied time stamp is a perfect match to the underlying cron
	// expression, and since this function is an entry point where `fromTime`
	// does not necessarily matches the underlying cron expression,
	// we first need to ensure supplied time stamp matches
	// the cron expression. If not, this means the supplied time
	// stamp falls in between matching time stamps, thus we move
	// to closest future matching immediately upon encountering a mismatching
	// time stamp.

	// year
	v := fromTime.Year()
	i := sort.SearchInts(expr.yearList, v)
	if i == len(expr.yearList) {
		return time.Time{}
	}
	if v != expr.yearList[i] {
		return expr.nextYear(fromTime)
	}
	// month
	v = int(fromTime.Month())
	i = sort.SearchInts(expr.monthList, v)
	if i == len(expr.monthList) {
		return expr.nextYear(fromTime)
	}
	if v != expr.monthList[i] {
		return expr.nextMonth(fromTime)
	}

	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(fromTime.Year(), int(fromTime.Month()))
	if len(expr.actualDaysOfMonthList) == 0 {
		return expr.nextMonth(fromTime)
	}

	// day of month
	v = fromTime.Day()
	i = sort.SearchInts(expr.actualDaysOfMonthList, v)
	if i == len(expr.actualDaysOfMonthList) {
		return expr.nextMonth(fromTime)
	}
	if v != expr.actualDaysOfMonthList[i] {
		return expr.nextDayOfMonth(fromTime)
	}
	// hour
	v = fromTime.Hour()
	i = sort.SearchInts(expr.hourList, v)
	if i == len(expr.hourList) {
		return expr.nextDayOfMonth(fromTime)
	}
	if v != expr.hourList[i] {
		return expr.nextHour(fromTime)
	}
	// minute
	v = fromTime.Minute()
	i = sort.SearchInts(expr.minuteList, v)
	if i == len(expr.minuteList) {
		return expr.nextHour(fromTime)
	}
	if v != expr.minuteList[i] {
		return expr.nextMinute(fromTime)
	}
	// second
	v = fromTime.Second()
	i = sort.SearchInts(expr.secondList, v)
	if i == len(expr.secondList) {
		return expr.nextMinute(fromTime)
	}

	// If we reach this point, there is nothing better to do
	// than to move to the next second

	return expr.nextSecond(fromTime)
}

/******************************************************************************/

func (expr *Expression) nextYear(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate year
	i := sort.SearchInts(expr.yearList, t.Year()+1)
	if i == len(expr.yearList) {
		return time.Time{}
	}
	// Year changed, need to recalculate actual days of month
	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(expr.yearList[i], expr.monthList[0])
	if len(expr.actualDaysOfMonthList) == 0 {
		return expr.nextMonth(time.Date(
			expr.yearList[i],
			time.Month(expr.monthList[0]),
			1,
			expr.hourList[0],
			expr.minuteList[0],
			expr.secondList[0],
			0,
			t.Location()))
	}
	return time.Date(
		expr.yearList[i],
		time.Month(expr.monthList[0]),
		expr.actualDaysOfMonthList[0],
		expr.hourList[0],
		expr.minuteList[0],
		expr.secondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextMonth(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate month
	i := sort.SearchInts(expr.monthList, int(t.Month())+1)
	if i == len(expr.monthList) {
		return expr.nextYear(t)
	}
	// Month changed, need to recalculate actual days of month
	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(t.Year(), expr.monthList[i])
	if len(expr.actualDaysOfMonthList) == 0 {
		return expr.nextMonth(time.Date(
			t.Year(),
			time.Month(expr.monthList[i]),
			1,
			expr.hourList[0],
			expr.minuteList[0],
			expr.secondList[0],
			0,
			t.Location()))
	}

	return time.Date(
		t.Year(),
		time.Month(expr.monthList[i]),
		expr.actualDaysOfMonthList[0],
		expr.hourList[0],
		expr.minuteList[0],
		expr.secondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextDayOfMonth(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate day of month
	i := sort.SearchInts(expr.actualDaysOfMonthList, t.Day()+1)
	if i == len(expr.actualDaysOfMonthList) {
		return expr.nextMonth(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		expr.actualDaysOfMonthList[i],
		expr.hourList[0],
		expr.minuteList[0],
		expr.secondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextHour(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate hour
	i := sort.SearchInts(expr.hourList, t.Hour()+1)
	if i == len(expr.hourList) {
		return expr.nextDayOfMonth(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		expr.hourList[i],
		expr.minuteList[0],
		expr.secondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextMinute(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate minute
	i := sort.SearchInts(expr.minuteList, t.Minute()+1)
	if i == len(expr.minuteList) {
		return expr.nextHour(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		expr.minuteList[i],
		expr.secondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextSecond(t time.Time) time.Time {
	// nextSecond() assumes all other fields are exactly matched
	// to the cron expression

	// Find index at which item in list is greater or equal to
	// candidate second
	i := sort.SearchInts(expr.secondList, t.Second()+1)
	if i == len(expr.secondList) {
		return expr.nextMinute(t)
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

/******************************************************************************/

func (expr *Expression) calculateActualDaysOfMonth(year, month int) []int {
	firstDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	allDays := genericDefaultList[1 : lastDayOfMonth.Day()+1]

	// A day must match both the day of month and the day of week, as every
	// other field must match. (The crontab man page runs the command when
	// either matches; netresearch/go-cron, by default, when both do.)

	// If both fields are not restricted, all days of the month are a hit
	if expr.daysOfMonthRestricted == false && expr.daysOfWeekRestricted == false {
		return allDays
	}
	if expr.daysOfMonthRestricted == false {
		return toList(expr.daysOfWeekInMonth(firstDayOfMonth, lastDayOfMonth))
	}
	daysOfMonth := expr.daysOfMonthInMonth(firstDayOfMonth, lastDayOfMonth)
	if expr.daysOfWeekRestricted == false {
		return toList(daysOfMonth)
	}
	daysOfWeek := expr.daysOfWeekInMonth(firstDayOfMonth, lastDayOfMonth)
	for v := range daysOfMonth {
		if !daysOfWeek[v] {
			delete(daysOfMonth, v)
		}
	}
	return toList(daysOfMonth)
}

// daysOfMonthInMonth returns the days of the month that the day-of-month
// field selects.
func (expr *Expression) daysOfMonthInMonth(firstDayOfMonth, lastDayOfMonth time.Time) map[int]bool {
	actualDaysOfMonthMap := make(map[int]bool)
	// Last day of month
	if expr.lastDayOfMonth {
		actualDaysOfMonthMap[lastDayOfMonth.Day()] = true
	}
	// Days before the last day of month: `L-3`
	for v := range expr.lastDayOfMonthOffsets {
		if v < lastDayOfMonth.Day() {
			actualDaysOfMonthMap[lastDayOfMonth.Day()-v] = true
		}
	}
	// Last work day of month
	if expr.lastWorkdayOfMonth {
		actualDaysOfMonthMap[workdayOfMonth(lastDayOfMonth, lastDayOfMonth)] = true
	}
	// Days of month
	for v := range expr.daysOfMonth {
		// Ignore days beyond end of month
		if v <= lastDayOfMonth.Day() {
			actualDaysOfMonthMap[v] = true
		}
	}
	// Work days of month
	// As per Wikipedia: month boundaries are not crossed.
	for v := range expr.workdaysOfMonth {
		// Ignore days beyond end of month
		if v <= lastDayOfMonth.Day() {
			actualDaysOfMonthMap[workdayOfMonth(firstDayOfMonth.AddDate(0, 0, v-1), lastDayOfMonth)] = true
		}
	}
	return actualDaysOfMonthMap
}

// daysOfWeekInMonth returns the days of the month that the day-of-week
// field selects.
func (expr *Expression) daysOfWeekInMonth(firstDayOfMonth, lastDayOfMonth time.Time) map[int]bool {
	actualDaysOfMonthMap := make(map[int]bool)
	// How far first sunday is from first day of month
	offset := 7 - int(firstDayOfMonth.Weekday())
	// days of week
	//  offset : (7 - day_of_week_of_1st_day_of_month)
	//  target : 1 + (7 * week_of_month) + (offset + day_of_week) % 7
	for v := range expr.daysOfWeek {
		w := dowNormalizedOffsets[(offset+v)%7]
		actualDaysOfMonthMap[w[0]] = true
		actualDaysOfMonthMap[w[1]] = true
		actualDaysOfMonthMap[w[2]] = true
		actualDaysOfMonthMap[w[3]] = true
		if len(w) > 4 && w[4] <= lastDayOfMonth.Day() {
			actualDaysOfMonthMap[w[4]] = true
		}
	}
	// days of week of specific week in the month
	//  offset : (7 - day_of_week_of_1st_day_of_month)
	//  target : 1 + (7 * week_of_month) + (offset + day_of_week) % 7
	for v := range expr.specificWeekDaysOfWeek {
		v = 1 + 7*(v/7) + (offset+v)%7
		if v <= lastDayOfMonth.Day() {
			actualDaysOfMonthMap[v] = true
		}
	}
	// Last days of week of the month
	lastWeekOrigin := firstDayOfMonth.AddDate(0, 1, -7)
	offset = 7 - int(lastWeekOrigin.Weekday())
	for v := range expr.lastWeekDaysOfWeek {
		v = lastWeekOrigin.Day() + (offset+v)%7
		if v <= lastDayOfMonth.Day() {
			actualDaysOfMonthMap[v] = true
		}
	}
	return actualDaysOfMonthMap
}

func workdayOfMonth(targetDom, lastDom time.Time) int {
	// If saturday, then friday
	// If sunday, then monday
	dom := targetDom.Day()
	dow := targetDom.Weekday()
	if dow == time.Saturday {
		if dom > 1 {
			dom -= 1
		} else {
			dom += 2
		}
	} else if dow == time.Sunday {
		if dom < lastDom.Day() {
			dom += 1
		} else {
			dom -= 2
		}
	}
	return dom
}
