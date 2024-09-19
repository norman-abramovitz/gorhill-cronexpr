NOTE
====

This cronexpr package is based on [Raymond Hill's][1] [cronexpr][2] package which
has been archived on Jan 10, 2019.  This package seems to be the closest Go package that attempts
to implement [Java's Quartz Cron expressions][3] which is similar to [Wikipedia's cron
 expanded definition][4].  The Go language [Quartz Scheduler][6], for example, does not come 
close to the [Java version][5]. As I evaluated other Go Language implementations, these are some
of the various features that were missing: 

1. seconds support
1. years support
1. oversized ranges
1. Last day (L) 
1. Last weekday (W)
1. Nearest Day (#)
1. Implementing AND/Intersection instead of OR/Union between day of month and day of the week fields.
As noted in the follwing [article][7] as is done in the [Cronitor Kubernetes implementation][8]

Guiding Principlese
===================

1. Try to stay true to the [Open Group][9], [Java][3] and [Wikipedia][4] specifications.
2. Added extensions to that keep to the spirit of those specifications.
3. Try to be responsive when people do PRs or Issue reporting  

[1]: https://github.com/gorhill
[2]: https://github.com/gorhill/cronexpr
[3]: https://www.quartz-scheduler.org/api/2.3.0/org/quartz/CronExpression.html
[4]: https://en.wikipedia.org/wiki/Cron#Cron_expression
[5]: https://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/crontrigger.html
[6]: https://github.com/reugn/go-quartz/?tab=readme-ov-file#cron-expression-format
[7]: https://crontab.guru/cron-bug.html
[8]: https://github.com/cronitorio/cronitor-kubernetes
[9]: https://pubs.opengroup.org/onlinepubs/9799919799/


The remaining package documentation will located the doc.go file.

Install
-------
    go get github.com/norman-abramovitz/cronexpr

Usage
-----
Import the library:

    import "github.com/gorhill/cronexpr"
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

The time zone of time values returned by `Next` and `NextN` is always the
time zone of the time value passed as argument, unless a zero time value is
returned.

API
---

TBD

License
-------

License: pick the one which suits you best:

- GPL v3 see <https://www.gnu.org/licenses/gpl.html>
- APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>

TODO
-------

1. Use Go Modules
2. Add 5 or 6 or 7 crontab fields
3. Add a cli command for quick testing
4. Add previous time support 
5. Add better timezone support
 
