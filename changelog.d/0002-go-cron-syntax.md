[Breaking Changes]
- The module is github.com/norman-abramovitz/gorhill-cronexpr/v2.
- Six fields read minute through year only when the sixth field is a
  year; otherwise they read seconds through day of the week.
- When day of month and day of week are both set, a day must match both
  (it matched either before).

[Features]
- The syntax of netresearch/go-cron's full parser, with its run times:
  wrapped ranges (50-10, FRI-MON), Sunday as 0 or 7, ?/n, L-n and #L,
  years up to 2147483647, @midnight, @every, @every_start,
  @every_aligned, @triggered, @manual and @none.
- A TZ= or CRON_TZ= prefix, whose zone the fields are matched in.
- H, H/n, H(a-b) and H(a-b)/n, from a hash key given to
  ParseWithHashKey; a key picks the same values as in go-cron.
- Prev and PrevN, the mirrors of Next and NextN.
- Daylight saving handled as go-cron does: a time repeated by a
  fall-back runs in both copies, and a time skipped by a spring-forward
  gap runs once per gap.

[BugFixes]
- The cronexpr command ran the upstream package instead of this one.
