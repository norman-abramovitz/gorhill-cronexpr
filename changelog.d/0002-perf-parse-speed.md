[Features]
- Parse is more than twice as fast (5.0 µs to 2.2 µs on the benchmark
  mix, 62 allocations to 19), and more than six times faster when parsing
  on many goroutines at once: it no longer takes a global lock (6.3 µs to
  0.9 µs per parse on 12 cores). Fields and entries are split without
  regular expressions; the layouts still match each entry. The layouts are compiled once at startup, at a
  cost of about 0.5 ms and 1 MB.
- Next is about 40% faster (491 ns to 278 ns) and no longer allocates.

[BugFixes]
- Next and Prev on one Expression shared between goroutines no longer
  return wrong runs (`0 0 L * *` gave the 30th of January). Upstream
  gorhill/cronexpr has the same race.
