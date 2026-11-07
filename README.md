Golang Cron expression parser
=============================
Given a cron expression and a time stamp, you can get the next time stamp which satisfies the cron expression.

In another project, I decided to use cron expression syntax to encode scheduling information. Thus this standalone library to parse and apply time stamps to cron expressions.

The time-matching algorithm in this implementation is efficient, it avoids as much as possible to guess the next matching time stamp, a common technique seen in a number of implementations out there.

There is also a companion command-line utility to evaluate cron time expressions: [cronexpr](cronexpr/) (which of course uses this library).

Implementation
--------------
The reference documentation for this implementation is found at
<https://en.wikipedia.org/wiki/Cron#CRON_expression>, which I copy/pasted here (laziness!) with modifications where this implementation differs:

    Field name     Mandatory?   Allowed values    Allowed special characters
    ----------     ----------   --------------    --------------------------
    Seconds        No           0-59              * ? / , - H
    Minutes        Yes          0-59              * ? / , - H
    Hours          Yes          0-23              * ? / , - H
    Day of month   Yes          1-31              * ? / , - L W H
    Month          Yes          1-12 or JAN-DEC   * ? / , - H
    Day of week    Yes          0-7 or SUN-SAT    * ? / , - L # H
    Year           No           1970–2099 for *   * ? / , -
                                1-2147483647 written

#### Asterisk ( * )
The asterisk indicates that the cron expression matches for all values of the field. E.g., using an asterisk in the 4th field (month) indicates every month. 

#### Question mark ( ? )
`?` is the same as `*`, and `?/n` the same as `*/n`.

#### Slash ( / )
Slashes describe increments of ranges. For example `3-59/15` in the minute field indicate the third minute of the hour and every 15 minutes thereafter. The form `*/...` is equivalent to the form "first-last/...", that is, an increment over the largest possible range of the field.

#### Comma ( , )
Commas are used to separate items of a list. For example, using `MON,WED,FRI` in the 5th field (day of week) means Mondays, Wednesdays and Fridays.

#### Hyphen ( - )
Hyphens define ranges. For example, 2000-2010 indicates every year between 2000 and 2010 AD, inclusive.

A range whose first value is above its last wraps past the end of the field: `50-10` in the minute field is minutes 50 to 59 and 0 to 10, `22-2` in the hour field runs over midnight, and `FRI-MON` is Friday to Monday. A step counts on across the wrap: `50-10/5` is minutes 50, 55, 0, 5 and 10. Years do not wrap.

#### L
`L` stands for "last". When used in the day-of-week field, it allows you to specify constructs such as "the last Friday" (`5L`, or `5#L`) of a given month. In the day-of-month field, it specifies the last day of the month, and `L-n` the day `n` days before it: `L-1` is the second-to-last day of the month.

#### W
The `W` character is allowed for the day-of-month field. This character is used to specify the business day (Monday-Friday) nearest the given day. As an example, if you were to specify `15W` as the value for the day-of-month field, the meaning is: "the nearest business day to the 15th of the month."

So, if the 15th is a Saturday, the trigger fires on Friday the 14th. If the 15th is a Sunday, the trigger fires on Monday the 16th. If the 15th is a Tuesday, then it fires on Tuesday the 15th. However if you specify `1W` as the value for day-of-month, and the 1st is a Saturday, the trigger fires on Monday the 3rd, as it does not 'jump' over the boundary of a month's days.

The `W` character can be specified only when the day-of-month is a single day, not a range or list of days.

The `W` character can also be combined with `L`, i.e. `LW` to mean "the last business day of the month."

#### Hash ( # )
`#` is allowed for the day-of-week field, and must be followed by a number between one and five, or by `L`. It allows you to specify constructs such as "the second Friday" of a given month.

#### H
`H` picks a value from a hash of a key, so that jobs given the same expression with different keys spread out instead of all running at once. It needs a key, given with `ParseWithHashKey`; `Parse` rejects an expression with `H`. The same key always picks the same values, and picks the ones that netresearch/go-cron picks.

* `H` is one value of the field: `H * * * *` runs once an hour, at a minute the key picks.
* `H/n` is every `n`-th value, from an offset the key picks: `H/15 * * * *` runs four times an hour.
* `H(a-b)` and `H(a-b)/n` do the same within `a` to `b`, which may wrap: `0 H(22-2) * * *` runs once a night.

`H` is allowed in every field but the year.

#### Day of month and day of week together
When both day-of-month and day-of-week are restricted, a day must match both, as every other field must match: `0 0 13 * FRI` runs on Friday the 13th. (The crontab of Vixie cron runs when either one matches.)

Predefined cron expressions
---------------------------
(Copied from <https://en.wikipedia.org/wiki/Cron#Predefined_scheduling_definitions>, with text modified according to this implementation) 

    Entry       Description                                                             Equivalent to
    @annually   Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *
    @yearly     Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *
    @monthly    Run once a month at midnight in the morning of the first of the month   0 0 0 1 * * *
    @weekly     Run once a week at midnight in the morning of Sunday                    0 0 0 * * 0 *
    @daily      Run once a day at midnight                                              0 0 0 * * * *
    @midnight   Run once a day at midnight                                              0 0 0 * * * *
    @hourly     Run once an hour at the beginning of the hour                           0 0 * * * * *
    @reboot     Not supported

These are not calendar expressions:

    Entry                Description
    @every d             Run every duration d after the time given, e.g. @every 1h30m
    @every_start d       The same as @every d
    @every_aligned d     Run at the multiples of d since 1970-01-01 00:00:00 UTC, e.g. @every_aligned 15m runs at :00, :15, :30 and :45
    @triggered           Never run by time: Next and Prev return the zero time. Also @manual and @none

The duration `d` is written as for Go's `time.ParseDuration`, must be at least one second, and is counted in whole seconds.

Time zones
----------
A `TZ=` or `CRON_TZ=` prefix gives the time zone whose wall clock the expression is matched against, by its IANA name, with or without quotes:

    TZ=America/New_York 0 9 * * MON-FRI
    CRON_TZ="Asia/Tokyo" @daily

Without a prefix, the expression is matched against the wall clock of the time zone of the time given.

Daylight saving time changes follow netresearch/go-cron:

* A time repeated when the clocks go back runs twice, once in each copy. `30 1 * * *` in New York runs at 01:30 EDT and again at 01:30 EST on the day summer time ends.
* A time skipped when the clocks go forward runs once, at the earliest skipped run moved on by the length of the gap. `30 2 * * *` in New York runs at 03:30 EDT on the day summer time starts.

Other details
-------------
* If only six fields are present, the sixth is read as a year when it is a number of 100 or more (also in a list, range or step: `2030-2035`). Then a `0` second field is prepended, that is, `* * * * * 2013` internally become `0 * * * * * 2013`. Otherwise the six fields are seconds through day of week, and a wildcard year field is appended, that is, `*/10 * * * * *` runs every ten seconds. A year below 100 needs all seven fields.
* If only five fields are present, a `0` second field is prepended and a wildcard year field is appended, that is, `* * * * Mon` internally become `0 * * * * Mon *`.
* Domain for day-of-week field is [0-7] instead of [0-6], 7 being Sunday (like 0). This to comply with http://linux.die.net/man/5/crontab#.
* A wildcard year runs from 1970 to 2099, but a year written out may be any year from 1 to 2147483647.
* As of now, the behavior of the code is undetermined if a malformed cron expression is supplied

Install
-------
    go get github.com/norman-abramovitz/gorhill-cronexpr/v2

Usage
-----
Import the library:

    import "github.com/norman-abramovitz/gorhill-cronexpr/v2"
    import "time"

Simplest way:

    nextTime := cronexpr.MustParse("0 0 29 2 *").Next(time.Now())

Assuming `time.Now()` is "2013-08-29 09:28:00", then `nextTime` will be "2016-02-29 00:00:00".

You can keep the returned Expression pointer around if you want to reuse it:

    expr := cronexpr.MustParse("0 0 29 2 *")
    nextTime := expr.Next(time.Now())
    ...
    nextTime = expr.Next(nextTime)

Use `time.IsZero()` to find out whether a valid time was returned. For example,

    cronexpr.MustParse("* * * * * 1980").Next(time.Now()).IsZero()

will return `true`, whereas

    cronexpr.MustParse("* * * * * 2050").Next(time.Now()).IsZero()

will return `false` (as of 2013-08-29...)

You may also query for `n` next time stamps:

    cronexpr.MustParse("0 0 29 2 *").NextN(time.Now(), 5)

which returns a slice of time.Time objects, containing the following time stamps (as of 2013-08-30):

    2016-02-29 00:00:00
    2020-02-29 00:00:00
    2024-02-29 00:00:00
    2028-02-29 00:00:00
    2032-02-29 00:00:00

`Prev` and `PrevN` are their mirrors: the closest time stamps before the time
given, latest first:

    cronexpr.MustParse("0 0 29 2 *").PrevN(time.Now(), 2)

The time zone of time values returned by `Next`, `NextN`, `Prev` and `PrevN`
is always the time zone of the time value passed as argument, unless a zero
time value is returned. This holds with a `TZ=` prefix too: the expression is
matched in the prefix's zone, and the run is returned in the zone of the time
given.

An expression with `H` is parsed with a key, such as the job's name:

    expr, err := cronexpr.ParseWithHashKey("H H(0-5) * * *", "nightly-backup")

API
---
<https://pkg.go.dev/github.com/norman-abramovitz/gorhill-cronexpr/v2>

License
-------

License: pick the one which suits you best:

- GPL v3 see <https://www.gnu.org/licenses/gpl.html>
- APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>

