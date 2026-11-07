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
