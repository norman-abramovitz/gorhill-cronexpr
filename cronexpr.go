/*!
 * Copyright 2013 Raymond Hill
 * Copyright 2026 Norman Abramovitz
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr.go
 * Version: 1.0
 * License: pick the one which suits you :
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

// Package cronexpr parses cron time expressions.
package cronexpr

/******************************************************************************/

import (
	"fmt"
	"time"
)

/******************************************************************************/

// A Expression represents a specific cron time expression as defined at
// <https://github.com/gorhill/cronexpr#implementation>
type Expression struct {
	expression             string
	kind                   int            // calendar, interval, aligned or never
	interval               time.Duration  // for interval and aligned
	location               *time.Location // from a TZ= prefix; nil: the zone of the time given
	hashValue              uint64         // for H, from the hash key
	hashKeyed              bool           // a hash key was given
	secondList             []int
	minuteList             []int
	hourList               []int
	daysOfMonth            map[int]bool
	workdaysOfMonth        map[int]bool
	lastDayOfMonth         bool
	lastDayOfMonthOffsets  map[int]bool
	lastWorkdayOfMonth     bool
	daysOfMonthRestricted  bool
	actualDaysOfMonthList  []int
	monthList              []int
	daysOfWeek             map[int]bool
	specificWeekDaysOfWeek map[int]bool
	lastWeekDaysOfWeek     map[int]bool
	daysOfWeekRestricted   bool
	yearList               []int
}

// Kinds of expression
const (
	calendar = iota // fields, or a descriptor such as @daily
	interval        // @every and @every_start: every interval from the time given
	aligned         // @every_aligned: the multiples of interval since the Unix epoch
	never           // @triggered, @manual, @none: never runs by time
)

/******************************************************************************/

// MustParse returns a new Expression pointer. It expects a well-formed cron
// expression. If a malformed cron expression is supplied, it will `panic`.
// See <https://github.com/gorhill/cronexpr#implementation> for documentation
// about what is a well-formed cron expression from this library's point of
// view.
func MustParse(cronLine string) *Expression {
	expr, err := Parse(cronLine)
	if err != nil {
		panic(err)
	}
	return expr
}

/******************************************************************************/

// Parse returns a new Expression pointer. An error is returned if a malformed
// cron expression is supplied.
// See <https://github.com/gorhill/cronexpr#implementation> for documentation
// about what is a well-formed cron expression from this library's point of
// view.
//
// An expression with H needs a hash key: use ParseWithHashKey.
func Parse(cronLine string) (*Expression, error) {
	return parse(cronLine, &Expression{})
}

// ParseWithHashKey is Parse for an expression that may hold H, which picks
// values from a hash of key: H, H/n, H(a-b) and H(a-b)/n, in any field but
// the year. The same key picks the same values as netresearch/go-cron.
func ParseWithHashKey(cronLine, key string) (*Expression, error) {
	if key == "" {
		return nil, fmt.Errorf("H needs a hash key that is not empty")
	}
	return parse(cronLine, &Expression{hashValue: hashOf(key), hashKeyed: true})
}

func parse(cronLine string, expr *Expression) (*Expression, error) {
	expr.expression = cronLine

	// Maybe a time zone prefix: TZ=America/New_York 0 9 * * *
	location, rest, err := parseTimeZone(cronLine)
	if err != nil {
		return nil, err
	}
	expr.location = location

	// Maybe a descriptor that is not a calendar
	if done, err := expr.parseIntervalDescriptor(rest); done {
		if err != nil {
			return nil, err
		}
		return expr, nil
	}

	// Maybe one of the built-in aliases is being used
	cron := cronNormalizer.Replace(rest)

	indices := fieldFinder.FindAllStringIndex(cron, -1)
	fieldCount := len(indices)
	if fieldCount < 5 {
		return nil, fmt.Errorf("missing field(s)")
	}
	// ignore fields beyond 7th
	if fieldCount > 7 {
		fieldCount = 7
	}

	var field = 0

	// second field (optional): with seven fields, or with six when the
	// sixth is not a year
	if fieldCount == 7 || (fieldCount == 6 && !looksLikeYear(cron[indices[5][0]:indices[5][1]])) {
		err = expr.secondFieldHandler(cron[indices[field][0]:indices[field][1]])
		if err != nil {
			return nil, err
		}
		field += 1
	} else {
		expr.secondList = []int{0}
	}

	// minute field
	err = expr.minuteFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// hour field
	err = expr.hourFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// day of month field
	err = expr.domFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// month field
	err = expr.monthFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// day of week field
	err = expr.dowFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// year field
	if field < fieldCount {
		err = expr.yearFieldHandler(cron[indices[field][0]:indices[field][1]])
		if err != nil {
			return nil, err
		}
	} else {
		expr.yearList = yearDescriptor.defaultList
	}

	return expr, nil
}

/******************************************************************************/

// Next returns the closest time instant immediately following `fromTime` which
// matches the cron expression `expr`.
//
// The `time.Location` of the returned time instant is the same as that of
// `fromTime`. The expression is matched against the wall clock of its TZ=
// zone, or of `fromTime`'s zone when it has none. A time skipped by a
// daylight saving spring-forward gap runs once per gap, shifted by the length
// of the gap; a time repeated by a fall-back runs in both copies.
//
// The zero value of time.Time is returned if no matching time instant exists
// or if a `fromTime` is itself a zero value.
func (expr *Expression) Next(fromTime time.Time) time.Time {
	// Special case
	if fromTime.IsZero() {
		return fromTime
	}
	switch expr.kind {
	case interval:
		return fromTime.Add(expr.interval - time.Duration(fromTime.Nanosecond()))
	case aligned:
		return fromTime.Add(expr.interval - gridOffset(fromTime, expr.interval))
	case never:
		return time.Time{}
	}

	t := fromTime.In(expr.zone(fromTime))
	next := expr.nextExisting(t)
	if run, ok := expr.nextGapRun(t, next); ok {
		next = run
	}
	if next.IsZero() {
		return next
	}
	return next.In(fromTime.Location())
}

// Prev returns the closest time instant before `fromTime`, and before the
// second it falls in, which matches the cron expression `expr`. It is the
// mirror of Next, and reports the same daylight saving runs.
//
// The zero value of time.Time is returned if no matching time instant exists
// or if a `fromTime` is itself a zero value.
func (expr *Expression) Prev(fromTime time.Time) time.Time {
	// Special case
	if fromTime.IsZero() {
		return fromTime
	}
	switch expr.kind {
	case interval:
		// NOTE: as go-cron's ConstantDelaySchedule.Prev, which adds the
		// fraction of a second back instead of dropping it
		return fromTime.Add(-expr.interval + time.Duration(fromTime.Nanosecond()))
	case aligned:
		off := gridOffset(fromTime, expr.interval)
		if off == 0 {
			off = expr.interval
		}
		return fromTime.Add(-off)
	case never:
		return time.Time{}
	}

	t := fromTime.In(expr.zone(fromTime))
	prev := expr.prevExisting(t)
	if run, ok := expr.prevGapRun(t, prev); ok {
		prev = run
	}
	if prev.IsZero() {
		return prev
	}
	return prev.In(fromTime.Location())
}

// zone returns the zone whose wall clock the expression is matched against.
func (expr *Expression) zone(fromTime time.Time) *time.Location {
	if expr.location != nil {
		return expr.location
	}
	return fromTime.Location()
}

/******************************************************************************/

// NextN returns a slice of `n` closest time instants immediately following
// `fromTime` which match the cron expression `expr`.
//
// The time instants in the returned slice are in chronological ascending order.
// The `time.Location` of the returned time instants is the same as that of
// `fromTime`.
//
// A slice with len between [0-`n`] is returned, that is, if not enough existing
// matching time instants exist, the number of returned entries will be less
// than `n`.
func (expr *Expression) NextN(fromTime time.Time, n uint) []time.Time {
	nextTimes := make([]time.Time, 0, n)
	for ; n > 0; n -= 1 {
		fromTime = expr.Next(fromTime)
		if fromTime.IsZero() {
			break
		}
		nextTimes = append(nextTimes, fromTime)
	}
	return nextTimes
}

// PrevN returns a slice of `n` closest time instants immediately preceding
// `fromTime` which match the cron expression `expr`. It is the mirror of
// NextN.
//
// The time instants in the returned slice are in chronological descending
// order. The `time.Location` of the returned time instants is the same as that
// of `fromTime`.
//
// A slice with len between [0-`n`] is returned, that is, if not enough existing
// matching time instants exist, the number of returned entries will be less
// than `n`.
func (expr *Expression) PrevN(fromTime time.Time, n uint) []time.Time {
	prevTimes := make([]time.Time, 0, n)
	for ; n > 0; n -= 1 {
		fromTime = expr.Prev(fromTime)
		if fromTime.IsZero() {
			break
		}
		prevTimes = append(prevTimes, fromTime)
	}
	return prevTimes
}
