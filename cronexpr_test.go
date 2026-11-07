/*!
 * Copyright 2013 Raymond Hill
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_test.go
 * Version: 1.0
 * License: pick the one which suits you best:
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

package cronexpr

/******************************************************************************/

import (
	"testing"
	"time"
)

/******************************************************************************/

type crontimes struct {
	from string
	next string
}

type crontest struct {
	expr   string
	layout string
	times  []crontimes
}

var crontests = []crontest{
	// Seconds
	{
		"* * * * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:00:01"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:59", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:59", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:59", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:59", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:59", "2013-01-01 00:00:00"},
		},
	},

	// every 5 Second
	{
		"*/5 * * * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:00:05"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:59", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:59", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:59", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:59", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:59", "2013-01-01 00:00:00"},
		},
	},

	// Minutes
	{
		"* * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:01:00"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:00", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:00", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:00", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:00", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:00", "2013-01-01 00:00:00"},
		},
	},

	// Minutes with interval
	{
		"17-43/5 * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:17:00"},
			{"2013-01-01 00:16:59", "2013-01-01 00:17:00"},
			{"2013-01-01 00:30:00", "2013-01-01 00:32:00"},
			{"2013-01-01 00:50:00", "2013-01-01 01:17:00"},
			{"2013-01-01 23:50:00", "2013-01-02 00:17:00"},
			{"2013-02-28 23:50:00", "2013-03-01 00:17:00"},
			{"2016-02-28 23:50:00", "2016-02-29 00:17:00"},
			{"2012-12-31 23:50:00", "2013-01-01 00:17:00"},
		},
	},

	// Minutes interval, list
	{
		"15-30/4,55 * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:15:00"},
			{"2013-01-01 00:16:00", "2013-01-01 00:19:00"},
			{"2013-01-01 00:30:00", "2013-01-01 00:55:00"},
			{"2013-01-01 00:55:00", "2013-01-01 01:15:00"},
			{"2013-01-01 23:55:00", "2013-01-02 00:15:00"},
			{"2013-02-28 23:55:00", "2013-03-01 00:15:00"},
			{"2016-02-28 23:55:00", "2016-02-29 00:15:00"},
			{"2012-12-31 23:54:00", "2012-12-31 23:55:00"},
			{"2012-12-31 23:55:00", "2013-01-01 00:15:00"},
		},
	},

	// Days of week
	{
		"0 0 * * MON",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Mon 2013-01-07 00:00"},
			{"2013-01-28 00:00:00", "Mon 2013-02-04 00:00"},
			{"2013-12-30 00:30:00", "Mon 2014-01-06 00:00"},
		},
	},
	{
		"0 0 * * friday",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Fri 2013-01-04 00:00"},
			{"2013-01-28 00:00:00", "Fri 2013-02-01 00:00"},
			{"2013-12-30 00:30:00", "Fri 2014-01-03 00:00"},
		},
	},
	{
		"0 0 * * 6,7",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Sat 2013-01-05 00:00"},
			{"2013-01-28 00:00:00", "Sat 2013-02-02 00:00"},
			{"2013-12-30 00:30:00", "Sat 2014-01-04 00:00"},
		},
	},

	// Specific days of week
	{
		"0 0 * * 6#5",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Sat 2013-11-30 00:00"},
		},
	},

	// Work day of month
	{
		"0 0 14W * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-03-31 00:00:00", "Mon 2013-04-15 00:00"},
			{"2013-08-31 00:00:00", "Fri 2013-09-13 00:00"},
		},
	},

	// Work day of month -- end of month
	{
		"0 0 30W * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-03-02 00:00:00", "Fri 2013-03-29 00:00"},
			{"2013-06-02 00:00:00", "Fri 2013-06-28 00:00"},
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2013-11-02 00:00:00", "Fri 2013-11-29 00:00"},
		},
	},

	// Last day of month
	{
		"0 0 L * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2014-01-01 00:00:00", "Fri 2014-01-31 00:00"},
			{"2014-02-01 00:00:00", "Fri 2014-02-28 00:00"},
			{"2016-02-15 00:00:00", "Mon 2016-02-29 00:00"},
		},
	},

	// Last work day of month
	{
		"0 0 LW * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2013-11-02 00:00:00", "Fri 2013-11-29 00:00"},
			{"2014-08-15 00:00:00", "Fri 2014-08-29 00:00"},
		},
	},

	// TODO: more tests
}

func TestExpressions(t *testing.T) {
	for _, test := range crontests {
		for _, times := range test.times {
			from, _ := time.Parse("2006-01-02 15:04:05", times.from)
			expr, err := Parse(test.expr)
			if err != nil {
				t.Errorf(`Parse("%s") returned "%s"`, test.expr, err.Error())
			}
			next := expr.Next(from)
			nextstr := next.Format(test.layout)
			if nextstr != times.next {
				t.Errorf(`("%s").Next("%s") = "%s", got "%s"`, test.expr, times.from, times.next, nextstr)
			}
		}
	}
}

/******************************************************************************/

func TestZero(t *testing.T) {
	from, _ := time.Parse("2006-01-02", "2013-08-31")
	next := MustParse("* * * * * 1980").Next(from)
	if next.IsZero() == false {
		t.Error(`("* * * * * 1980").Next("2013-08-31").IsZero() returned 'false', expected 'true'`)
	}

	next = MustParse("* * * * * 2050").Next(from)
	if next.IsZero() == true {
		t.Error(`("* * * * * 2050").Next("2013-08-31").IsZero() returned 'true', expected 'false'`)
	}

	next = MustParse("* * * * * 2099").Next(time.Time{})
	if next.IsZero() == false {
		t.Error(`("* * * * * 2014").Next(time.Time{}).IsZero() returned 'true', expected 'false'`)
	}
}

/******************************************************************************/

func TestNextN(t *testing.T) {
	expected := []string{
		"Sat, 30 Nov 2013 00:00:00",
		"Sat, 29 Mar 2014 00:00:00",
		"Sat, 31 May 2014 00:00:00",
		"Sat, 30 Aug 2014 00:00:00",
		"Sat, 29 Nov 2014 00:00:00",
	}
	from, _ := time.Parse("2006-01-02 15:04:05", "2013-09-02 08:44:30")
	result := MustParse("0 0 * * 6#5").NextN(from, uint(len(expected)))
	if len(result) != len(expected) {
		t.Errorf(`MustParse("0 0 * * 6#5").NextN("2013-09-02 08:44:30", 5):\n"`)
		t.Errorf(`  Expected %d returned time values but got %d instead`, len(expected), len(result))
	}
	for i, next := range result {
		nextStr := next.Format("Mon, 2 Jan 2006 15:04:15")
		if nextStr != expected[i] {
			t.Errorf(`MustParse("0 0 * * 6#5").NextN("2013-09-02 08:44:30", 5):\n"`)
			t.Errorf(`  result[%d]: expected "%s" but got "%s"`, i, expected[i], nextStr)
		}
	}
}

func TestNextN_every5min(t *testing.T) {
	expected := []string{
		"Mon, 2 Sep 2013 08:45:00",
		"Mon, 2 Sep 2013 08:50:00",
		"Mon, 2 Sep 2013 08:55:00",
		"Mon, 2 Sep 2013 09:00:00",
		"Mon, 2 Sep 2013 09:05:00",
	}
	from, _ := time.Parse("2006-01-02 15:04:05", "2013-09-02 08:44:32")
	result := MustParse("*/5 * * * *").NextN(from, uint(len(expected)))
	if len(result) != len(expected) {
		t.Errorf(`MustParse("*/5 * * * *").NextN("2013-09-02 08:44:30", 5):\n"`)
		t.Errorf(`  Expected %d returned time values but got %d instead`, len(expected), len(result))
	}
	for i, next := range result {
		nextStr := next.Format("Mon, 2 Jan 2006 15:04:05")
		if nextStr != expected[i] {
			t.Errorf(`MustParse("*/5 * * * *").NextN("2013-09-02 08:44:30", 5):\n"`)
			t.Errorf(`  result[%d]: expected "%s" but got "%s"`, i, expected[i], nextStr)
		}
	}
}

// Issue: https://github.com/gorhill/cronexpr/issues/16
func TestInterval_Interval60Issue(t *testing.T) {
	_, err := Parse("*/60 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}

	_, err = Parse("*/61 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 61 should return err")
	}

	_, err = Parse("2/60 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}

	_, err = Parse("2-20/61 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}
}

/******************************************************************************/

var benchmarkExpressions = []string{
	"* * * * *",
	"@hourly",
	"@weekly",
	"@yearly",
	"30 3 15W 3/3 *",
	"30 0 0 1-31/5 Oct-Dec * 2000,2006,2008,2013-2015",
	"0 0 0 * Feb-Nov/2 thu#3 2000-2050",
}
var benchmarkExpressionsLen = len(benchmarkExpressions)

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MustParse(benchmarkExpressions[i%benchmarkExpressionsLen])
	}
}

func BenchmarkNext(b *testing.B) {
	exprs := make([]*Expression, benchmarkExpressionsLen)
	for i := 0; i < benchmarkExpressionsLen; i++ {
		exprs[i] = MustParse(benchmarkExpressions[i])
	}
	from := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := exprs[i%benchmarkExpressionsLen]
		next := expr.Next(from)
		next = expr.Next(next)
		next = expr.Next(next)
		next = expr.Next(next)
		next = expr.Next(next)
	}
}

/******************************************************************************/

// The syntax and run times of netresearch/go-cron's full parser. The expected
// times are go-cron's, except where noted.

type runtest struct {
	spec string
	key  string // hash key for H, "" for Parse
	zone string
	from string // in zone, 2006-01-02 15:04:05
	prev bool   // Prev instead of Next
	want []string
}

const runLayout = "2006-01-02 15:04:05 MST"

var runtests = []runtest{
	// Six fields: a year in the sixth makes it minute through year,
	// otherwise seconds through day of the week
	{"*/10 * * * * *", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 00:00:10 UTC", "2026-01-01 00:00:20 UTC"}},
	{"0 0 0 * * *", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-03 00:00:00 UTC"}},
	{"0 0 1 1 * 2030", "", "UTC", "2026-01-01 00:00:00", false, []string{"2030-01-01 00:00:00 UTC", "zero"}},
	{"0 0 1 1 * */2", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 01:00:00 UTC", "2026-02-01 01:00:00 UTC"}},

	// Day of the month and day of the week both set: both must match
	{"0 0 13 * 5", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-02-13 00:00:00 UTC", "2026-03-13 00:00:00 UTC"}},
	{"0 0 1-7 * 1", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-05 00:00:00 UTC", "2026-02-02 00:00:00 UTC"}},
	{"0 0 L * 5", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-07-31 00:00:00 UTC", "2027-04-30 00:00:00 UTC"}},

	// Wrapped ranges
	{"50-10 * * * *", "", "UTC", "2026-01-01 00:08:00", false, []string{"2026-01-01 00:09:00 UTC", "2026-01-01 00:10:00 UTC", "2026-01-01 00:50:00 UTC", "2026-01-01 00:51:00 UTC"}},
	{"50-10/5 * * * *", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 00:05:00 UTC", "2026-01-01 00:10:00 UTC", "2026-01-01 00:50:00 UTC", "2026-01-01 00:55:00 UTC", "2026-01-01 01:00:00 UTC"}},
	{"0 22-2 * * *", "", "UTC", "2026-01-01 03:00:00", false, []string{"2026-01-01 22:00:00 UTC", "2026-01-01 23:00:00 UTC", "2026-01-02 00:00:00 UTC", "2026-01-02 01:00:00 UTC", "2026-01-02 02:00:00 UTC"}},
	{"0 0 28-3 * *", "", "UTC", "2026-01-04 00:00:00", false, []string{"2026-01-28 00:00:00 UTC", "2026-01-29 00:00:00 UTC", "2026-01-30 00:00:00 UTC", "2026-01-31 00:00:00 UTC", "2026-02-01 00:00:00 UTC", "2026-02-02 00:00:00 UTC", "2026-02-03 00:00:00 UTC", "2026-02-28 00:00:00 UTC"}},
	{"0 0 1 11-2 *", "", "UTC", "2026-03-01 00:00:00", false, []string{"2026-11-01 00:00:00 UTC", "2026-12-01 00:00:00 UTC", "2027-01-01 00:00:00 UTC", "2027-02-01 00:00:00 UTC"}},
	{"0 0 * * FRI-MON", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-03 00:00:00 UTC", "2026-01-04 00:00:00 UTC", "2026-01-05 00:00:00 UTC", "2026-01-09 00:00:00 UTC"}},

	// Sunday as 0 and 7
	{"0 0 * * 0-7", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-03 00:00:00 UTC", "2026-01-04 00:00:00 UTC"}},
	{"0 0 * * 5-7", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-03 00:00:00 UTC", "2026-01-04 00:00:00 UTC", "2026-01-09 00:00:00 UTC"}},
	{"0 0 * * 6-1", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-03 00:00:00 UTC", "2026-01-04 00:00:00 UTC", "2026-01-05 00:00:00 UTC", "2026-01-10 00:00:00 UTC"}},
	{"0 0 * * 1/2", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-04 00:00:00 UTC", "2026-01-05 00:00:00 UTC", "2026-01-07 00:00:00 UTC", "2026-01-09 00:00:00 UTC"}},
	{"0 0 * * */2", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-03 00:00:00 UTC", "2026-01-04 00:00:00 UTC", "2026-01-06 00:00:00 UTC", "2026-01-08 00:00:00 UTC", "2026-01-10 00:00:00 UTC"}},
	{"0 0 * * 5-1/2", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-04 00:00:00 UTC", "2026-01-09 00:00:00 UTC", "2026-01-11 00:00:00 UTC"}},
	{"0 0 * * 7-2/2", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-04 00:00:00 UTC", "2026-01-06 00:00:00 UTC", "2026-01-11 00:00:00 UTC", "2026-01-13 00:00:00 UTC"}},

	// ?/n, L-n, #L
	{"0 0 ?/2 * *", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-03 00:00:00 UTC", "2026-01-05 00:00:00 UTC", "2026-01-07 00:00:00 UTC"}},
	{"0 0 L-3 * *", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-28 00:00:00 UTC", "2026-02-25 00:00:00 UTC"}},
	{"0 0 * * 5#L", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-30 00:00:00 UTC", "2026-02-27 00:00:00 UTC"}},
	{"0 0 * * FRI#L", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-30 00:00:00 UTC"}},

	// Years outside 1970-2099. go-cron's Next stops 5 years ahead and
	// returns zero for 2100 and 2500; cronexpr runs them.
	{"0 0 1 1 * 2100", "", "UTC", "2026-01-01 00:00:00", false, []string{"2100-01-01 00:00:00 UTC"}},
	{"0 0 1 1 * 2500", "", "UTC", "2026-01-01 00:00:00", false, []string{"2500-01-01 00:00:00 UTC"}},
	{"0 0 1 1 * 1969", "", "UTC", "2026-01-01 00:00:00", false, []string{"zero"}},

	// Descriptors
	{"@midnight", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-03 00:00:00 UTC"}},
	{"@every 1h", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 01:00:00 UTC", "2026-01-01 02:00:00 UTC"}},
	{"@every 90m", "", "UTC", "2026-01-01 00:20:00", false, []string{"2026-01-01 01:50:00 UTC", "2026-01-01 03:20:00 UTC"}},
	{"@every_start 1h", "", "UTC", "2026-01-01 00:20:00", false, []string{"2026-01-01 01:20:00 UTC"}},
	{"@every_aligned 1h", "", "UTC", "2026-01-01 00:20:00", false, []string{"2026-01-01 01:00:00 UTC", "2026-01-01 02:00:00 UTC"}},
	{"@triggered", "", "UTC", "2026-01-01 00:00:00", false, []string{"zero"}},
	{"@every 1h", "", "UTC", "2026-01-01 01:00:00", true, []string{"2026-01-01 00:00:00 UTC", "2025-12-31 23:00:00 UTC"}},
	{"@every_aligned 1h", "", "UTC", "2026-01-01 01:20:00", true, []string{"2026-01-01 01:00:00 UTC", "2026-01-01 00:00:00 UTC"}},

	// Time zone prefix: the run is returned in the zone of the time given
	{"TZ=America/New_York 0 9 * * *", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 14:00:00 UTC", "2026-01-02 14:00:00 UTC"}},
	{"CRON_TZ='Asia/Tokyo' 0 9 * * *", "", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC"}},

	// H, with a hash key
	{"H * * * *", "job-a", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 00:54:00 UTC", "2026-01-01 01:54:00 UTC"}},
	{"H H(0-5) * * *", "job-a", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 00:54:00 UTC", "2026-01-02 00:54:00 UTC"}},
	{"H/15 * * * *", "job-b", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 00:08:00 UTC", "2026-01-01 00:23:00 UTC", "2026-01-01 00:38:00 UTC", "2026-01-01 00:53:00 UTC"}},
	{"0 0 H(1-28) * *", "job-c", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-05 00:00:00 UTC", "2026-02-05 00:00:00 UTC"}},
	{"0 0 * * H", "job-d", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-09 00:00:00 UTC"}},
	{"0 0 * * H/2", "job-d", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-02 00:00:00 UTC", "2026-01-05 00:00:00 UTC", "2026-01-07 00:00:00 UTC", "2026-01-09 00:00:00 UTC"}},
	{"0 H(22-2) * * *", "job-e", "UTC", "2026-01-01 00:00:00", false, []string{"2026-01-01 01:00:00 UTC", "2026-01-02 01:00:00 UTC"}},

	// DST: a time skipped by a spring-forward gap runs once, shifted by the
	// gap; a time repeated by a fall-back runs in both copies
	{"30 2 * * *", "", "America/New_York", "2026-03-08 00:00:00", false, []string{"2026-03-08 03:30:00 EDT", "2026-03-09 02:30:00 EDT", "2026-03-10 02:30:00 EDT"}},
	{"10,45 2,3 * * *", "", "America/New_York", "2026-03-08 01:00:00", false, []string{"2026-03-08 03:10:00 EDT", "2026-03-08 03:45:00 EDT", "2026-03-09 02:10:00 EDT", "2026-03-09 02:45:00 EDT"}},
	{"20,40 2 * * *", "", "Australia/Lord_Howe", "2026-10-04 00:00:00", false, []string{"2026-10-04 02:40:00 +11", "2026-10-04 02:50:00 +11", "2026-10-05 02:20:00 +11", "2026-10-05 02:40:00 +11"}},
	{"30 1 * * *", "", "America/New_York", "2026-11-01 00:00:00", false, []string{"2026-11-01 01:30:00 EDT", "2026-11-01 01:30:00 EST", "2026-11-02 01:30:00 EST"}},
	{"*/15 * * * *", "", "America/New_York", "2026-11-01 01:40:00", false, []string{"2026-11-01 01:45:00 EDT", "2026-11-01 01:00:00 EST", "2026-11-01 01:15:00 EST", "2026-11-01 01:30:00 EST", "2026-11-01 01:45:00 EST", "2026-11-01 02:00:00 EST", "2026-11-01 02:15:00 EST"}},
	{"30 2 * * *", "", "America/New_York", "2026-03-09 00:00:00", true, []string{"2026-03-08 03:30:00 EDT", "2026-03-07 02:30:00 EST"}},
	{"30 1 * * *", "", "America/New_York", "2026-11-02 00:00:00", true, []string{"2026-11-01 01:30:00 EST", "2026-11-01 01:30:00 EDT", "2026-10-31 01:30:00 EDT"}},
	{"*/15 * * * *", "", "America/New_York", "2026-11-01 01:20:00", true, []string{"2026-11-01 01:15:00 EDT", "2026-11-01 01:00:00 EDT", "2026-11-01 00:45:00 EDT", "2026-11-01 00:30:00 EDT", "2026-11-01 00:15:00 EDT", "2026-11-01 00:00:00 EDT", "2026-10-31 23:45:00 EDT"}},
	{"20,40 2 * * *", "", "Australia/Lord_Howe", "2026-10-05 00:00:00", true, []string{"2026-10-04 02:50:00 +11", "2026-10-04 02:40:00 +11", "2026-10-03 02:40:00 +1030", "2026-10-03 02:20:00 +1030"}},

	// Prev
	{"0 0 * * MON", "", "UTC", "2013-01-08 00:00:00", true, []string{"2013-01-07 00:00:00 UTC", "2012-12-31 00:00:00 UTC"}},
	{"17-43/5 * * * *", "", "UTC", "2013-01-01 01:00:00", true, []string{"2013-01-01 00:42:00 UTC", "2013-01-01 00:37:00 UTC", "2013-01-01 00:32:00 UTC"}},
	{"0 0 L * *", "", "UTC", "2014-03-01 00:00:00", true, []string{"2014-02-28 00:00:00 UTC", "2014-01-31 00:00:00 UTC", "2013-12-31 00:00:00 UTC"}},
	{"0 0 * * 6#5", "", "UTC", "2014-01-01 00:00:00", true, []string{"2013-11-30 00:00:00 UTC", "2013-08-31 00:00:00 UTC"}},
	{"0 0 14W * *", "", "UTC", "2013-05-01 00:00:00", true, []string{"2013-04-15 00:00:00 UTC", "2013-03-14 00:00:00 UTC"}},
	{"* * * * * 2050", "", "UTC", "2013-08-31 00:00:00", true, []string{"zero"}},
}

func TestRuns(t *testing.T) {
	for _, test := range runtests {
		loc, err := time.LoadLocation(test.zone)
		if err != nil {
			t.Fatal(err)
		}
		var expr *Expression
		if test.key != "" {
			expr, err = ParseWithHashKey(test.spec, test.key)
		} else {
			expr, err = Parse(test.spec)
		}
		if err != nil {
			t.Errorf("Parse(%q) returned %q", test.spec, err)
			continue
		}
		at, _ := time.ParseInLocation("2006-01-02 15:04:05", test.from, loc)
		for i, want := range test.want {
			if test.prev {
				at = expr.Prev(at)
			} else {
				at = expr.Next(at)
			}
			got := "zero"
			if !at.IsZero() {
				got = at.Format(runLayout)
			}
			if got != want {
				t.Errorf("%q from %s, run %d (prev=%v): want %s, got %s", test.spec, test.from, i+1, test.prev, want, got)
				break
			}
			if at.IsZero() {
				break
			}
		}
	}
}

func TestRejected(t *testing.T) {
	for _, spec := range []string{
		"@every 500ms",        // below one second
		"H * * * *",           // H without a hash key
		"0 0 1 1 * 2030-2020", // year ranges do not wrap
		"0 0 L-31 * *",
		"TZ=Nowhere/Special 0 9 * * *",
		"TZ=America/New_York",
	} {
		if _, err := Parse(spec); err == nil {
			t.Errorf("Parse(%q) returned no error", spec)
		}
	}
	if _, err := ParseWithHashKey("H * * * *", ""); err == nil {
		t.Errorf(`ParseWithHashKey("H * * * *", "") returned no error`)
	}
}

// A run's Prev and Next lead back to it, for every spec of crontests.
func TestPrevUndoesNext(t *testing.T) {
	for _, test := range crontests {
		expr := MustParse(test.expr)
		for _, times := range test.times {
			from, _ := time.Parse("2006-01-02 15:04:05", times.from)
			run := expr.Next(from)
			if back := expr.Next(expr.Prev(run)); !back.Equal(run) {
				t.Errorf("%q: Next(Prev(%s)) = %s", test.expr, run, back)
			}
		}
	}
}
