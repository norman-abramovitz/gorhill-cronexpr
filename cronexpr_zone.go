/*!
 * Copyright 2026 Norman Abramovitz
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_zone.go
 * License: pick the one which suits you :
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

package cronexpr

/******************************************************************************/

import (
	"fmt"
	"math/bits"
	"regexp"
	"strings"
	"time"
)

// Time zones and daylight saving, with the rules of netresearch/go-cron:
//
//   - The expression matches the wall clock. Within one zone period (a span
//     with one UTC offset) the wall clock is the instant plus that offset, so
//     nextWall and prevWall find the run, and the offset turns it back into an
//     instant. A run that falls outside the period is searched again in the
//     next (or previous) one.
//   - A wall time repeated by a fall-back exists twice, and runs twice.
//   - A wall time skipped by a spring-forward gap does not exist. The gap
//     runs once instead: at its earliest skipped run, shifted by the length
//     of the gap (02:30 runs at 03:30 after a one-hour gap at 02:00).

/******************************************************************************/

// The end of the wildcard year list, where searches stop
var searchEnd = time.Date(yearDescriptor.max+1, 1, 1, 0, 0, 0, 0, time.UTC)

// wall returns the wall clock of t, to the second, carried in UTC.
func wall(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}

// offsetOf returns the UTC offset in effect at t, in seconds.
func offsetOf(t time.Time) int {
	_, offset := t.Zone()
	return offset
}

// instant returns the time whose wall clock is w under the UTC offset.
func instant(w time.Time, offset int, loc *time.Location) time.Time {
	return time.Unix(w.Unix()-int64(offset), 0).In(loc)
}

// zoneEndAfter returns the first zone change after t, or the zero time when
// none follows. Past a zone's stored changes, ZoneBounds reports bounds at
// year ends and, just after one, an end that is not after t; probe past it.
func zoneEndAfter(t time.Time) time.Time {
	probe := t
	for i := 0; i < 48; i++ {
		_, end := probe.ZoneBounds()
		if end.IsZero() || end.After(t) {
			return end
		}
		probe = probe.Add(time.Hour)
	}
	return time.Time{}
}

/******************************************************************************/

// nextExisting returns the first instant after t whose wall clock matches.
func (expr *Expression) nextExisting(t time.Time) time.Time {
	after := wall(t)
	for {
		w := expr.nextWall(after)
		if w.IsZero() {
			return w
		}
		run := instant(w, offsetOf(t), t.Location())
		end := zoneEndAfter(t)
		if end.IsZero() || run.Before(end) {
			return run
		}
		// Past this zone period: search the next one from its first second
		t = end
		after = wall(end).Add(-time.Second)
	}
}

// prevExisting returns the last instant before t's second whose wall clock
// matches.
func (expr *Expression) prevExisting(t time.Time) time.Time {
	before := wall(t)
	for {
		w := expr.prevWall(before)
		if w.IsZero() {
			return w
		}
		run := instant(w, offsetOf(t), t.Location())
		start, _ := t.ZoneBounds()
		if start.IsZero() || !run.Before(start) {
			return run
		}
		// Before this zone period: search the previous one from its end
		t = start.Add(-time.Nanosecond)
		before = wall(t).Add(time.Second)
	}
}

/******************************************************************************/

// gapLength returns how far the zone change at tr moves the clock forward:
// the length of a spring-forward gap, or zero or less for any other change.
func gapLength(tr time.Time) time.Duration {
	return time.Duration(offsetOf(tr)-offsetOf(tr.Add(-time.Nanosecond))) * time.Second
}

// gapRun returns the run for the spring-forward gap that opens at the zone
// change tr: the earliest wall time the gap skips that matches, shifted by
// the length of the gap. It reports false when tr is no gap, or the gap
// skips no run.
func (expr *Expression) gapRun(tr time.Time) (time.Time, bool) {
	gap := gapLength(tr)
	if gap <= 0 {
		return time.Time{}, false
	}
	skippedFrom := wall(tr.Add(-time.Nanosecond)).Add(time.Second)
	w := expr.nextWall(skippedFrom.Add(-time.Second))
	if w.IsZero() || !w.Before(skippedFrom.Add(gap)) {
		return time.Time{}, false
	}
	return tr.Add(w.Sub(skippedFrom)), true
}

// nextGapRun returns the earliest gap run after t and before next; a zero
// next means no run follows t.
func (expr *Expression) nextGapRun(t, next time.Time) (time.Time, bool) {
	limit := next
	if limit.IsZero() {
		limit = searchEnd
	}
	// The run for the gap that opened t's zone period can lie after t (for
	// 01:59:59 -> 03:00, a 02:30 run is at 03:30).
	tr, _ := t.ZoneBounds()
	if tr.IsZero() || tr.After(t) || !t.Before(tr.Add(gapLength(tr))) {
		tr = zoneEndAfter(t)
	}
	for !tr.IsZero() && tr.Before(limit) {
		if run, ok := expr.gapRun(tr); ok && run.After(t) && run.Before(limit) {
			return run, true
		}
		tr = zoneEndAfter(tr)
	}
	return time.Time{}, false
}

// prevGapRun returns the latest gap run before t's second and after prev; a
// zero prev means no run precedes t.
func (expr *Expression) prevGapRun(t, prev time.Time) (time.Time, bool) {
	limit := prev
	if limit.IsZero() {
		limit = time.Date(yearDescriptor.min, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	before := t.Truncate(time.Second)
	tr, _ := t.ZoneBounds()
	for !tr.IsZero() {
		if run, ok := expr.gapRun(tr); ok && run.Before(before) && run.After(limit) {
			return run, true
		}
		if !tr.After(limit) {
			break
		}
		tr, _ = tr.Add(-time.Nanosecond).ZoneBounds()
	}
	return time.Time{}, false
}

/******************************************************************************/

var (
	// `TZ=America/New_York 0 9 * * *`, `CRON_TZ="Asia/Tokyo" @daily`
	timeZonePrefix = regexp.MustCompile(`^(?:CRON_)?TZ=(\S*) (.*)$`)
	// IANA names: letters, digits, / _ - + :
	timeZoneName = regexp.MustCompile(`^[A-Za-z0-9/_+:-]{1,64}$`)
)

// parseTimeZone splits a TZ= or CRON_TZ= prefix from the expression. Matching
// quotes around the name are dropped. It returns a nil location when there is
// no prefix.
func parseTimeZone(cronLine string) (*time.Location, string, error) {
	if !timeZoneStart.MatchString(cronLine) {
		return nil, cronLine, nil
	}
	parts := timeZonePrefix.FindStringSubmatch(cronLine)
	if parts == nil || strings.TrimSpace(parts[2]) == "" {
		return nil, "", fmt.Errorf("missing fields after time zone: '%s'", cronLine)
	}
	name := parts[1]
	if len(name) >= 2 && (name[0] == '"' || name[0] == '\'') && name[len(name)-1] == name[0] {
		name = name[1 : len(name)-1]
	}
	if !timeZoneName.MatchString(name) {
		return nil, "", fmt.Errorf("invalid time zone '%s'", name)
	}
	location, err := time.LoadLocation(name)
	if err != nil {
		return nil, "", fmt.Errorf("unknown time zone '%s'", name)
	}
	return location, strings.TrimSpace(parts[2]), nil
}

var timeZoneStart = regexp.MustCompile(`^(?:CRON_)?TZ=`)

/******************************************************************************/

// gridOffset returns how far t lies past the last multiple of d since the
// Unix epoch, in [0, d), for d > 0. It works on seconds and nanoseconds
// apart, so it holds outside the int64 nanosecond range too.
func gridOffset(t time.Time, d time.Duration) time.Duration {
	n := uint64(d)
	secRem := t.Unix() % int64(d)
	if secRem < 0 {
		secRem += int64(d)
	}
	hi, lo := bits.Mul64(uint64(secRem), uint64(time.Second)%n)
	rem := bits.Rem64(hi, lo, n)
	rem = (rem + uint64(t.Nanosecond())%n) % n
	return time.Duration(rem)
}
