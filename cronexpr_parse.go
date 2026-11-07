/*!
 * Copyright 2013 Raymond Hill
 * Copyright 2026 Norman Abramovitz
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_parse.go
 * Version: 1.0
 * License: pick the one which suits you best:
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

package cronexpr

/******************************************************************************/

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

/******************************************************************************/

var (
	genericDefaultList = []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
	}
	yearDefaultList = []int{
		1970, 1971, 1972, 1973, 1974, 1975, 1976, 1977, 1978, 1979,
		1980, 1981, 1982, 1983, 1984, 1985, 1986, 1987, 1988, 1989,
		1990, 1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998, 1999,
		2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009,
		2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019,
		2020, 2021, 2022, 2023, 2024, 2025, 2026, 2027, 2028, 2029,
		2030, 2031, 2032, 2033, 2034, 2035, 2036, 2037, 2038, 2039,
		2040, 2041, 2042, 2043, 2044, 2045, 2046, 2047, 2048, 2049,
		2050, 2051, 2052, 2053, 2054, 2055, 2056, 2057, 2058, 2059,
		2060, 2061, 2062, 2063, 2064, 2065, 2066, 2067, 2068, 2069,
		2070, 2071, 2072, 2073, 2074, 2075, 2076, 2077, 2078, 2079,
		2080, 2081, 2082, 2083, 2084, 2085, 2086, 2087, 2088, 2089,
		2090, 2091, 2092, 2093, 2094, 2095, 2096, 2097, 2098, 2099,
	}
)

/******************************************************************************/

var (
	numberTokens = map[string]int{
		"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"00": 0, "01": 1, "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9,
		"10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16, "17": 17, "18": 18, "19": 19,
		"20": 20, "21": 21, "22": 22, "23": 23, "24": 24, "25": 25, "26": 26, "27": 27, "28": 28, "29": 29,
		"30": 30, "31": 31, "32": 32, "33": 33, "34": 34, "35": 35, "36": 36, "37": 37, "38": 38, "39": 39,
		"40": 40, "41": 41, "42": 42, "43": 43, "44": 44, "45": 45, "46": 46, "47": 47, "48": 48, "49": 49,
		"50": 50, "51": 51, "52": 52, "53": 53, "54": 54, "55": 55, "56": 56, "57": 57, "58": 58, "59": 59,
		"1970": 1970, "1971": 1971, "1972": 1972, "1973": 1973, "1974": 1974, "1975": 1975, "1976": 1976, "1977": 1977, "1978": 1978, "1979": 1979,
		"1980": 1980, "1981": 1981, "1982": 1982, "1983": 1983, "1984": 1984, "1985": 1985, "1986": 1986, "1987": 1987, "1988": 1988, "1989": 1989,
		"1990": 1990, "1991": 1991, "1992": 1992, "1993": 1993, "1994": 1994, "1995": 1995, "1996": 1996, "1997": 1997, "1998": 1998, "1999": 1999,
		"2000": 2000, "2001": 2001, "2002": 2002, "2003": 2003, "2004": 2004, "2005": 2005, "2006": 2006, "2007": 2007, "2008": 2008, "2009": 2009,
		"2010": 2010, "2011": 2011, "2012": 2012, "2013": 2013, "2014": 2014, "2015": 2015, "2016": 2016, "2017": 2017, "2018": 2018, "2019": 2019,
		"2020": 2020, "2021": 2021, "2022": 2022, "2023": 2023, "2024": 2024, "2025": 2025, "2026": 2026, "2027": 2027, "2028": 2028, "2029": 2029,
		"2030": 2030, "2031": 2031, "2032": 2032, "2033": 2033, "2034": 2034, "2035": 2035, "2036": 2036, "2037": 2037, "2038": 2038, "2039": 2039,
		"2040": 2040, "2041": 2041, "2042": 2042, "2043": 2043, "2044": 2044, "2045": 2045, "2046": 2046, "2047": 2047, "2048": 2048, "2049": 2049,
		"2050": 2050, "2051": 2051, "2052": 2052, "2053": 2053, "2054": 2054, "2055": 2055, "2056": 2056, "2057": 2057, "2058": 2058, "2059": 2059,
		"2060": 2060, "2061": 2061, "2062": 2062, "2063": 2063, "2064": 2064, "2065": 2065, "2066": 2066, "2067": 2067, "2068": 2068, "2069": 2069,
		"2070": 2070, "2071": 2071, "2072": 2072, "2073": 2073, "2074": 2074, "2075": 2075, "2076": 2076, "2077": 2077, "2078": 2078, "2079": 2079,
		"2080": 2080, "2081": 2081, "2082": 2082, "2083": 2083, "2084": 2084, "2085": 2085, "2086": 2086, "2087": 2087, "2088": 2088, "2089": 2089,
		"2090": 2090, "2091": 2091, "2092": 2092, "2093": 2093, "2094": 2094, "2095": 2095, "2096": 2096, "2097": 2097, "2098": 2098, "2099": 2099,
	}
	monthTokens = map[string]int{
		`1`: 1, `jan`: 1, `january`: 1,
		`2`: 2, `feb`: 2, `february`: 2,
		`3`: 3, `mar`: 3, `march`: 3,
		`4`: 4, `apr`: 4, `april`: 4,
		`5`: 5, `may`: 5,
		`6`: 6, `jun`: 6, `june`: 6,
		`7`: 7, `jul`: 7, `july`: 7,
		`8`: 8, `aug`: 8, `august`: 8,
		`9`: 9, `sep`: 9, `september`: 9,
		`10`: 10, `oct`: 10, `october`: 10,
		`11`: 11, `nov`: 11, `november`: 11,
		`12`: 12, `dec`: 12, `december`: 12,
	}
	dowTokens = map[string]int{
		`0`: 0, `sun`: 0, `sunday`: 0,
		`1`: 1, `mon`: 1, `monday`: 1,
		`2`: 2, `tue`: 2, `tuesday`: 2,
		`3`: 3, `wed`: 3, `wednesday`: 3,
		`4`: 4, `thu`: 4, `thursday`: 4,
		`5`: 5, `fri`: 5, `friday`: 5,
		`6`: 6, `sat`: 6, `saturday`: 6,
		`7`: 7, `07`: 7,
		`00`: 0, `01`: 1, `02`: 2, `03`: 3, `04`: 4, `05`: 5, `06`: 6,
	}
)

/******************************************************************************/

func atoi(s string) int {
	return numberTokens[s]
}

type fieldDescriptor struct {
	name         string
	min, max     int
	defaultList  []int
	valuePattern string
	atoi         func(string) int
	cyclic       bool // a range may wrap past the maximum: `22-2`
	layouts      *fieldLayouts
}

var (
	secondDescriptor = fieldDescriptor{
		name:         "second",
		min:          0,
		max:          59,
		defaultList:  genericDefaultList[0:60],
		valuePattern: `0?[0-9]|[1-5][0-9]`,
		atoi:         atoi,
		cyclic:       true,
	}
	minuteDescriptor = fieldDescriptor{
		name:         "minute",
		min:          0,
		max:          59,
		defaultList:  genericDefaultList[0:60],
		valuePattern: `0?[0-9]|[1-5][0-9]`,
		atoi:         atoi,
		cyclic:       true,
	}
	hourDescriptor = fieldDescriptor{
		name:         "hour",
		min:          0,
		max:          23,
		defaultList:  genericDefaultList[0:24],
		valuePattern: `0?[0-9]|1[0-9]|2[0-3]`,
		atoi:         atoi,
		cyclic:       true,
	}
	domDescriptor = fieldDescriptor{
		name:         "day-of-month",
		min:          1,
		max:          31,
		defaultList:  genericDefaultList[1:32],
		valuePattern: `0?[1-9]|[12][0-9]|3[01]`,
		atoi:         atoi,
		cyclic:       true,
	}
	monthDescriptor = fieldDescriptor{
		name:         "month",
		min:          1,
		max:          12,
		defaultList:  genericDefaultList[1:13],
		valuePattern: `0?[1-9]|1[012]|jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|january|february|march|april|march|april|june|july|august|september|october|november|december`,
		atoi: func(s string) int {
			return monthTokens[s]
		},
		cyclic: true,
	}
	// 0 and 7 are both Sunday; 7 becomes 0 once the values are built
	dowDescriptor = fieldDescriptor{
		name:         "day-of-week",
		min:          0,
		max:          7,
		defaultList:  genericDefaultList[0:7],
		valuePattern: `0?[0-7]|sun|mon|tue|wed|thu|fri|sat|sunday|monday|tuesday|wednesday|thursday|friday|saturday`,
		atoi: func(s string) int {
			return dowTokens[s]
		},
		cyclic: true,
	}
	// A wildcard year is 1970 to 2099; a written year is any year from 1 to
	// 2147483647, as netresearch/go-cron
	yearDescriptor = fieldDescriptor{
		name:         "year",
		min:          1970,
		max:          2099,
		defaultList:  yearDefaultList[:],
		valuePattern: `[0-9]+`,
		atoi: func(s string) int {
			v, err := strconv.Atoi(s)
			if err != nil {
				return -1
			}
			return v
		},
	}
)

/******************************************************************************/

var (
	layoutWildcard             = `^\*$|^\?$`
	layoutValue                = `^(%value%)$`
	layoutRange                = `^(%value%)-(%value%)$`
	layoutWildcardAndInterval  = `^[*?]/(\d+)$`
	layoutValueAndInterval     = `^(%value%)/(\d+)$`
	layoutRangeAndInterval     = `^(%value%)-(%value%)/(\d+)$`
	layoutLastDom              = `^l$`
	layoutLastDomOffset        = `^l-(0?[1-9]|[12][0-9]|30)$`
	layoutWorkdom              = `^(%value%)w$`
	layoutLastWorkdom          = `^lw$`
	layoutDowOfLastWeek        = `^(%value%)(?:l|#l)$`
	layoutDowOfSpecificWeek    = `^(%value%)#([1-5])$`
	layoutHash                 = `^h$`
	layoutHashAndInterval      = `^h/(\d+)$`
	layoutHashRange            = `^h\((%value%)-(%value%)\)$`
	layoutHashRangeAndInterval = `^h\((%value%)-(%value%)\)/(\d+)$`
)

// isFieldSpace reports whether r separates fields: the regexp `\s` set of
// space, tab, newline, carriage return and form feed. strings.Fields would
// also split on a vertical tab and Unicode spaces.
func isFieldSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f'
}

/******************************************************************************/

var cronNormalizer = strings.NewReplacer(
	"@yearly", "0 0 0 1 1 * *",
	"@annually", "0 0 0 1 1 * *",
	"@monthly", "0 0 0 1 * * *",
	"@weekly", "0 0 0 * * 0 *",
	"@daily", "0 0 0 * * * *",
	"@midnight", "0 0 0 * * * *",
	"@hourly", "0 0 * * * * *")

/******************************************************************************/

var (
	// `@every 1h30m`, `@every_start 90s`, `@every_aligned 15m`
	layoutIntervalDescriptor = regexp.MustCompile(`^@(every|every_start|every_aligned) (.*)$`)
	// Run only when started by hand
	layoutNeverDescriptor = regexp.MustCompile(`^@(triggered|manual|none)$`)
)

// parseIntervalDescriptor reads a descriptor that is not a calendar. It
// reports false when s is none of them.
func (expr *Expression) parseIntervalDescriptor(s string) (bool, error) {
	if layoutNeverDescriptor.MatchString(s) {
		expr.kind = never
		return true, nil
	}
	parts := layoutIntervalDescriptor.FindStringSubmatch(s)
	if parts == nil {
		return false, nil
	}
	d, err := time.ParseDuration(parts[2])
	if err != nil {
		return true, fmt.Errorf("invalid duration in '%s'", s)
	}
	// At least one second, counted in whole seconds
	if d < time.Second {
		return true, fmt.Errorf("@%s duration must be at least 1s: '%s'", parts[1], s)
	}
	expr.interval = d - d%time.Second
	expr.kind = interval
	if parts[1] == "every_aligned" {
		expr.kind = aligned
	}
	return true, nil
}

/******************************************************************************/

// looksLikeYear reports whether a field holds a year: a number of 100 or
// more. It tells a sixth field that is a year from a day of week.
func looksLikeYear(field string) bool {
	for _, part := range yearFinder.FindAllString(field, -1) {
		if v, err := strconv.Atoi(part); err == nil && v >= 100 {
			return true
		}
	}
	return false
}

var yearFinder = regexp.MustCompile(`[0-9]+`)

/******************************************************************************/

func (expr *Expression) secondFieldHandler(s string) error {
	var err error
	expr.secondList, err = expr.genericFieldHandler(s, secondDescriptor)
	return err
}

/******************************************************************************/

func (expr *Expression) minuteFieldHandler(s string) error {
	var err error
	expr.minuteList, err = expr.genericFieldHandler(s, minuteDescriptor)
	return err
}

/******************************************************************************/

func (expr *Expression) hourFieldHandler(s string) error {
	var err error
	expr.hourList, err = expr.genericFieldHandler(s, hourDescriptor)
	return err
}

/******************************************************************************/

func (expr *Expression) monthFieldHandler(s string) error {
	var err error
	expr.monthList, err = expr.genericFieldHandler(s, monthDescriptor)
	return err
}

/******************************************************************************/

// The most years a year field may select, as netresearch/go-cron
const maxYears = 100000

func (expr *Expression) yearFieldHandler(s string) error {
	directives, err := genericFieldParse(s, &yearDescriptor, nil)
	if err != nil {
		return err
	}
	values := make(map[int]bool) // a written year reaches 2147483647: no valueSet
	for _, directive := range directives {
		switch directive.kind {
		case one:
			directive.last = directive.first
			directive.step = 1
		case span:
		case all:
			expr.yearList = yearDescriptor.defaultList
			return nil
		default:
			return fmt.Errorf("syntax error in year field: '%s'", s[directive.sbeg:directive.send])
		}
		// A written year is any year from 1 to 2147483647; year ranges do
		// not wrap
		if directive.first < 1 || directive.last > math.MaxInt32 {
			return fmt.Errorf("year out of range in year field: '%s'", s[directive.sbeg:directive.send])
		}
		if directive.first > directive.last {
			return fmt.Errorf("year range ends before it starts: '%s'", s[directive.sbeg:directive.send])
		}
		if (directive.last-directive.first)/directive.step+1+len(values) > maxYears {
			return fmt.Errorf("year field selects more than %d years: '%s'", maxYears, s)
		}
		for year := directive.first; year <= directive.last; year += directive.step {
			values[year] = true
		}
	}
	expr.yearList = toList(values)
	return nil
}

/******************************************************************************/

const (
	none   = 0
	one    = 1
	span   = 2
	all    = 3
	picked = 4 // values picked by H
	// day of month only
	lastDom       = 5 // `L`
	lastDomOffset = 6 // `L-3`: first is the offset
	lastWorkdom   = 7 // `LW`
	workdom       = 8 // `15W`: first is the day
	// day of week only
	dowOfLastWeek     = 9  // `5L`: first is the day
	dowOfSpecificWeek = 10 // `5#3`: first is week*7 + day, less 7
)

type cronDirective struct {
	kind   int
	first  int
	last   int
	step   int
	values []int
	sbeg   int
	send   int
}

func (expr *Expression) genericFieldHandler(s string, desc fieldDescriptor) ([]int, error) {
	directives, err := genericFieldParse(s, &desc, expr)
	if err != nil {
		return nil, err
	}
	var values valueSet
	for _, directive := range directives {
		switch directive.kind {
		case one:
			populateOne(&values, directive.first)
		case span:
			populateSpan(&values, directive.first, directive.last, directive.step, desc)
		case all:
			return desc.defaultList, nil
		case picked:
			populateList(&values, directive.values)
		default:
			return nil, fmt.Errorf("syntax error in %s field: '%s'", desc.name, s[directive.sbeg:directive.send])
		}
	}
	return values.list(), nil
}

func (expr *Expression) dowFieldHandler(s string) error {
	expr.daysOfWeekRestricted = true
	expr.daysOfWeek = valueSet{}
	expr.lastWeekDaysOfWeek = valueSet{}
	expr.specificWeekDaysOfWeek = valueSet{}

	directives, err := genericFieldParse(s, &dowDescriptor, expr)
	if err != nil {
		return err
	}
	for _, directive := range directives {
		switch directive.kind {
		case dowOfLastWeek:
			populateOne(&expr.lastWeekDaysOfWeek, directive.first)
		case dowOfSpecificWeek:
			populateOne(&expr.specificWeekDaysOfWeek, directive.first)
		case one:
			populateOne(&expr.daysOfWeek, directive.first)
		case span:
			populateSpan(&expr.daysOfWeek, directive.first, directive.last, directive.step, dowDescriptor)
		case all:
			populateMany(&expr.daysOfWeek, directive.first, directive.last, directive.step)
			expr.daysOfWeekRestricted = false
		case picked:
			populateList(&expr.daysOfWeek, directive.values)
		default:
			return fmt.Errorf("syntax error in day-of-week field: '%s'", s[directive.sbeg:directive.send])
		}
	}
	// Sunday is 0 and 7: fold 7 into 0 once the values are built
	if expr.daysOfWeek[7] {
		expr.daysOfWeek[7] = false
		expr.daysOfWeek[0] = true
	}
	return nil
}

func (expr *Expression) domFieldHandler(s string) error {
	expr.daysOfMonthRestricted = true
	expr.lastDayOfMonth = false
	expr.lastWorkdayOfMonth = false
	expr.daysOfMonth = valueSet{}           // days of month
	expr.workdaysOfMonth = valueSet{}       // work days of month
	expr.lastDayOfMonthOffsets = valueSet{} // days before the last day of month

	directives, err := genericFieldParse(s, &domDescriptor, expr)
	if err != nil {
		return err
	}

	for _, directive := range directives {
		switch directive.kind {
		case lastDom:
			expr.lastDayOfMonth = true
		case lastDomOffset:
			populateOne(&expr.lastDayOfMonthOffsets, directive.first)
		case lastWorkdom:
			expr.lastWorkdayOfMonth = true
		case workdom:
			populateOne(&expr.workdaysOfMonth, directive.first)
		case one:
			populateOne(&expr.daysOfMonth, directive.first)
		case span:
			populateSpan(&expr.daysOfMonth, directive.first, directive.last, directive.step, domDescriptor)
		case all:
			populateMany(&expr.daysOfMonth, directive.first, directive.last, directive.step)
			expr.daysOfMonthRestricted = false
		case picked:
			populateList(&expr.daysOfMonth, directive.values)
		default:
			return fmt.Errorf("syntax error in day-of-month field: '%s'", s[directive.sbeg:directive.send])
		}
	}
	return nil
}

/******************************************************************************/

// valueSet marks the values a field selects: set[v] is true when v is
// selected. It holds every field but the year: seconds and minutes go to
// 59, and `5#3` is kept as week*7 + day of week, at most 34.
type valueSet [60]bool

// list returns the selected values, lowest first.
func (set *valueSet) list() []int {
	count := 0
	for _, selected := range set {
		if selected {
			count++
		}
	}
	list := make([]int, 0, count)
	for v, selected := range set {
		if selected {
			list = append(list, v)
		}
	}
	return list
}

func populateOne(values *valueSet, v int) {
	values[v] = true
}

func populateMany(values *valueSet, min, max, step int) {
	for i := min; i <= max; i += step {
		values[i] = true
	}
}

func populateList(values *valueSet, list []int) {
	for _, v := range list {
		values[v] = true
	}
}

// populateSpan adds the values from first to last by step. When first is
// past last, the range wraps around: `22-2` is 22, 23, 0, 1, 2, and the step
// counts on across the wrap. A day of week step walks the seven distinct
// days, so Sunday (0 and 7) counts once.
func populateSpan(values *valueSet, first, last, step int, desc fieldDescriptor) {
	if first <= last {
		populateMany(values, first, last, step)
		return
	}
	if step > 1 {
		first, last, desc = dowCycle(first, last, desc)
		if first <= last {
			populateMany(values, first, last, step)
			return
		}
	}
	for v := first; v <= last+desc.max-desc.min+1; v += step {
		values[wrapValue(v, desc)] = true
	}
}

func toList(set map[int]bool) []int {
	list := make([]int, len(set))
	i := 0
	for k := range set {
		list[i] = k
		i += 1
	}
	sort.Ints(list)
	return list
}

/******************************************************************************/

// The forms an entry of a field can take, as a recognizer reports them.
const (
	formNone              = iota // no form: a syntax error
	formWildcard                 // `*`, `?`
	formValue                    // `5`
	formRange                    // `5-20`, or `22-2` wrapping around
	formWildcardStep             // `*/2`
	formValueStep                // `5/2`
	formRangeStep                // `5-20/2`
	formHash                     // `H`
	formHashStep                 // `H/15`
	formHashRange                // `H(0-29)`
	formHashRangeStep            // `H(0-29)/10`
	formLastDom                  // `L`
	formLastDomOffset            // `L-3`
	formLastWorkdom              // `LW`
	formWorkdom                  // `15W`
	formDowOfLastWeek            // `5L`, `5#L`
	formDowOfSpecificWeek        // `5#3`: last is the week
)

// fieldEntry is one entry of a field as a recognizer found it: its form, the
// text of its values and step, and its place in the field.
type fieldEntry struct {
	form              int
	first, last, step string
	sbeg, send        int
}

// genericFieldParse reads the entries of a field. recognizeField finds each
// entry's form; this turns the text it found into values. expr gives the
// hash for H; it is nil for the year field, where H is a syntax error.
func genericFieldParse(s string, desc *fieldDescriptor, expr *Expression) ([]*cronDirective, error) {
	var buf [8]fieldEntry
	entries, err := recognizeField(s, desc, buf[:0])
	if err != nil {
		return nil, err
	}
	directives := make([]*cronDirective, 0, len(entries))
	for i := range entries {
		entry := &entries[i]
		directive := cronDirective{
			sbeg: entry.sbeg,
			send: entry.send,
		}
		step := func() error {
			directive.step = atoi(entry.step)
			if directive.step < 1 || directive.step > desc.max {
				return fmt.Errorf("invalid interval %s", s[entry.sbeg:entry.send])
			}
			return nil
		}
		switch entry.form {
		case formWildcard:
			directive.kind = all
			directive.first = desc.min
			directive.last = desc.max
			directive.step = 1
		case formValue:
			directive.kind = one
			directive.first = desc.atoi(entry.first)
		case formRange:
			directive.kind = span
			directive.first = desc.atoi(entry.first)
			directive.last = desc.atoi(entry.last)
			directive.step = 1
		case formWildcardStep:
			directive.kind = span
			directive.first = desc.min
			directive.last = desc.max
			if err := step(); err != nil {
				return nil, err
			}
		case formValueStep:
			directive.kind = span
			directive.first = desc.atoi(entry.first)
			directive.last = desc.max
			if err := step(); err != nil {
				return nil, err
			}
		case formRangeStep:
			directive.kind = span
			directive.first = desc.atoi(entry.first)
			directive.last = desc.atoi(entry.last)
			if err := step(); err != nil {
				return nil, err
			}
		case formHash, formHashStep, formHashRange, formHashRangeStep:
			if expr != nil {
				if err := expr.parseHash(s[entry.sbeg:entry.send], entry, desc, &directive); err != nil {
					return nil, err
				}
			}
		case formLastDom:
			directive.kind = lastDom
		case formLastDomOffset:
			directive.kind = lastDomOffset
			directive.first = atoi(entry.first)
		case formLastWorkdom:
			directive.kind = lastWorkdom
		case formWorkdom:
			directive.kind = workdom
			directive.first = desc.atoi(entry.first)
		case formDowOfLastWeek:
			directive.kind = dowOfLastWeek
			directive.first = desc.atoi(entry.first)
		case formDowOfSpecificWeek:
			directive.kind = dowOfSpecificWeek
			directive.first = (desc.atoi(entry.last)-1)*7 + desc.atoi(entry.first)%7
		}
		// formNone: no behavior for this one, let caller deal with it
		directives = append(directives, &directive)
	}
	// At least one entry must be present
	if len(directives) == 0 {
		return nil, fmt.Errorf("%s field: missing directive", desc.name)
	}
	return directives, nil
}

// parseHash reads an H entry into directive.
func (expr *Expression) parseHash(sentry string, entry *fieldEntry, desc *fieldDescriptor, directive *cronDirective) error {
	first, last, step := desc.min, desc.max, 0
	switch entry.form {
	case formHashStep:
		step, _ = strconv.Atoi(entry.step)
	case formHashRange:
		first = desc.atoi(entry.first)
		last = desc.atoi(entry.last)
	case formHashRangeStep:
		first = desc.atoi(entry.first)
		last = desc.atoi(entry.last)
		step, _ = strconv.Atoi(entry.step)
	}
	if !expr.hashKeyed {
		return fmt.Errorf("H needs a hash key, use ParseWithHashKey: '%s'", sentry)
	}
	directive.kind = picked
	switch {
	case entry.step != "" && step < 1:
		return fmt.Errorf("invalid interval %s", sentry)
	case step > 0:
		directive.values = hashStep(first, last, step, *desc, expr.hashValue)
	default:
		directive.values = []int{hashOne(first, last, *desc, expr.hashValue)}
	}
	return nil
}

/******************************************************************************/

// regexRecognizeField finds the form of each entry of a field by trying the
// field's layouts in turn. Entries are split on commas; empty ones (`1,,2`)
// are skipped. sbeg and send keep each entry's place in s for error messages.
func regexRecognizeField(s string, layouts *fieldLayouts, entries []fieldEntry) []fieldEntry {
	send := -1
	for sentry := range strings.SplitSeq(s, ",") {
		sbeg := send + 1
		send = sbeg + len(sentry)
		if sentry == "" {
			continue
		}
		entry := fieldEntry{sbeg: sbeg, send: send}
		snormal := strings.ToLower(sentry)
		// group returns the text of capture group i
		var pairs []int
		group := func(i int) string { return snormal[pairs[2*i]:pairs[2*i+1]] }
		match := func(layout *regexp.Regexp) bool {
			pairs = layout.FindStringSubmatchIndex(snormal)
			return pairs != nil
		}
		switch {
		case layouts.wildcard.MatchString(snormal):
			entry.form = formWildcard
		case layouts.value.MatchString(snormal):
			entry.form = formValue
			entry.first = snormal
		case match(layouts.valueRange):
			entry.form = formRange
			entry.first, entry.last = group(1), group(2)
		case match(layouts.wildcardAndInterval):
			entry.form = formWildcardStep
			entry.step = group(1)
		case match(layouts.valueAndInterval):
			entry.form = formValueStep
			entry.first, entry.step = group(1), group(2)
		case match(layouts.rangeAndInterval):
			entry.form = formRangeStep
			entry.first, entry.last, entry.step = group(1), group(2), group(3)
		case layouts.hash.MatchString(snormal):
			entry.form = formHash
		case match(layouts.hashAndInterval):
			entry.form = formHashStep
			entry.step = group(1)
		case match(layouts.hashRange):
			entry.form = formHashRange
			entry.first, entry.last = group(1), group(2)
		case match(layouts.hashRangeAndInterval):
			entry.form = formHashRangeStep
			entry.first, entry.last, entry.step = group(1), group(2), group(3)
		case layouts.lastDom.MatchString(snormal):
			entry.form = formLastDom
		case match(layouts.lastDomOffset):
			entry.form = formLastDomOffset
			entry.first = group(1)
		case layouts.lastWorkdom.MatchString(snormal):
			entry.form = formLastWorkdom
		case match(layouts.workdom):
			entry.form = formWorkdom
			entry.first = group(1)
		case match(layouts.dowOfLastWeek):
			entry.form = formDowOfLastWeek
			entry.first = group(1)
		case match(layouts.dowOfSpecificWeek):
			entry.form = formDowOfSpecificWeek
			entry.first, entry.last = group(1), group(2)
		}
		entries = append(entries, entry)
	}
	return entries
}

/******************************************************************************/

// fieldLayouts holds the layouts compiled for one field, with the field's
// value pattern in place of %value%.
type fieldLayouts struct {
	wildcard             *regexp.Regexp // layoutWildcard
	value                *regexp.Regexp // layoutValue
	valueRange           *regexp.Regexp // layoutRange
	wildcardAndInterval  *regexp.Regexp // layoutWildcardAndInterval
	valueAndInterval     *regexp.Regexp // layoutValueAndInterval
	rangeAndInterval     *regexp.Regexp // layoutRangeAndInterval
	lastDom              *regexp.Regexp // layoutLastDom
	lastDomOffset        *regexp.Regexp // layoutLastDomOffset
	workdom              *regexp.Regexp // layoutWorkdom
	lastWorkdom          *regexp.Regexp // layoutLastWorkdom
	dowOfLastWeek        *regexp.Regexp // layoutDowOfLastWeek
	dowOfSpecificWeek    *regexp.Regexp // layoutDowOfSpecificWeek
	hash                 *regexp.Regexp // layoutHash
	hashAndInterval      *regexp.Regexp // layoutHashAndInterval
	hashRange            *regexp.Regexp // layoutHashRange
	hashRangeAndInterval *regexp.Regexp // layoutHashRangeAndInterval
}

// compileLayouts compiles a field's layouts. compiled, when not nil, shares
// the patterns several fields have in common.
func compileLayouts(desc *fieldDescriptor, compiled map[string]*regexp.Regexp) *fieldLayouts {
	compile := func(layout string) *regexp.Regexp {
		pattern := strings.Replace(layout, `%value%`, desc.valuePattern, -1)
		if compiled == nil {
			return regexp.MustCompile(pattern)
		}
		if compiled[pattern] == nil {
			compiled[pattern] = regexp.MustCompile(pattern)
		}
		return compiled[pattern]
	}
	return &fieldLayouts{
		wildcard:             compile(layoutWildcard),
		value:                compile(layoutValue),
		valueRange:           compile(layoutRange),
		wildcardAndInterval:  compile(layoutWildcardAndInterval),
		valueAndInterval:     compile(layoutValueAndInterval),
		rangeAndInterval:     compile(layoutRangeAndInterval),
		lastDom:              compile(layoutLastDom),
		lastDomOffset:        compile(layoutLastDomOffset),
		workdom:              compile(layoutWorkdom),
		lastWorkdom:          compile(layoutLastWorkdom),
		dowOfLastWeek:        compile(layoutDowOfLastWeek),
		dowOfSpecificWeek:    compile(layoutDowOfSpecificWeek),
		hash:                 compile(layoutHash),
		hashAndInterval:      compile(layoutHashAndInterval),
		hashRange:            compile(layoutHashRange),
		hashRangeAndInterval: compile(layoutHashRangeAndInterval),
	}
}

// The layouts are compiled once at startup and only read afterwards, so
// parsing needs no lock. A field reaches them without a lookup.
func init() {
	compiled := make(map[string]*regexp.Regexp) // fields share most patterns
	for _, desc := range fieldDescriptors {
		desc.layouts = compileLayouts(desc, compiled)
	}
}

var fieldDescriptors = []*fieldDescriptor{
	&secondDescriptor, &minuteDescriptor, &hourDescriptor, &domDescriptor,
	&monthDescriptor, &dowDescriptor, &yearDescriptor,
}
