# Benchmarks and Test Cases

It seems that Go's otherwise excellent
[testing](https://golang.org/pkg/testing/)
infrastructure is susceptible to certain "contamination effects"
in which the presence (or absence) of *test cases* influences the
performance measured by *benchmarks*.
Check it out:

```
$ benchstat without.txt with.txt
name               old time/op    new time/op    delta
PushFrontQueue-2     85.5µs ± 1%    97.7µs ± 1%  +14.31%  (p=0.000 n=10+10)
PushFrontList-2       159µs ± 0%     167µs ± 1%   +4.93%  (p=0.000 n=9+10)
PushBackQueue-2      85.9µs ± 1%    98.3µs ± 1%  +14.38%  (p=0.000 n=10+10)
PushBackList-2        159µs ± 0%     166µs ± 1%   +4.91%  (p=0.000 n=9+9)
PushBackChannel-2     117µs ± 1%     146µs ± 2%  +25.01%  (p=0.000 n=10+10)
RandomQueue-2         153µs ± 0%     174µs ± 3%  +13.69%  (p=0.000 n=9+9)
RandomList-2          284µs ± 1%     297µs ± 1%   +4.68%  (p=0.000 n=10+8)

name               old alloc/op   new alloc/op   delta
PushFrontQueue-2     40.9kB ± 0%    40.9kB ± 0%     ~     (all equal)
PushFrontList-2      57.4kB ± 0%    57.4kB ± 0%     ~     (all equal)
PushBackQueue-2      40.9kB ± 0%    40.9kB ± 0%     ~     (all equal)
PushBackList-2       57.4kB ± 0%    57.4kB ± 0%     ~     (all equal)
PushBackChannel-2    24.7kB ± 0%    24.7kB ± 0%     ~     (all equal)
RandomQueue-2        45.7kB ± 0%    45.7kB ± 0%     ~     (all equal)
RandomList-2         90.8kB ± 0%    90.8kB ± 0%     ~     (all equal)

name               old allocs/op  new allocs/op  delta
PushFrontQueue-2      1.03k ± 0%     1.03k ± 0%     ~     (all equal)
PushFrontList-2       2.05k ± 0%     2.05k ± 0%     ~     (all equal)
PushBackQueue-2       1.03k ± 0%     1.03k ± 0%     ~     (all equal)
PushBackList-2        2.05k ± 0%     2.05k ± 0%     ~     (all equal)
PushBackChannel-2     1.03k ± 0%     1.03k ± 0%     ~     (all equal)
RandomQueue-2         1.63k ± 0%     1.63k ± 0%     ~     (all equal)
RandomList-2          3.24k ± 0%     3.24k ± 0%     ~     (all equal)
```

The *only* difference between `old time` and `new time` is that the
latter had two additional *test cases* in `queue_test.go`.
The presence of those test cases makes *all* benchmarks appear worse
for some reason.
If I comment them out and run the *exact* *same* benchmarks again,
everything is "back to normal" as it were.
For now I am reporting the *worse* results in the official `README` just
to be conservative, but I believe something needs to be fixed on the Go
side of things.
