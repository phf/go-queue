# Benchmarks and Test Cases

It seems that *some* versions of Go's otherwise excellent
[testing](https://golang.org/pkg/testing/)
infrastructure are susceptible to certain "contamination effects"
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
$ go version
go version go1.7.5 linux/amd64
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

**Update**: Something *was* apparently fixed in a more recent version
of Go. Check it out:

```
$ benchstat without.txt with.txt
name               old time/op    new time/op    delta
PushFrontQueue-2     82.7µs ± 1%    82.8µs ± 1%    ~     (p=0.832 n=19+16)
PushFrontList-2       161µs ± 1%     162µs ± 1%  +0.48%  (p=0.003 n=19+18)
PushBackQueue-2      83.3µs ± 1%    83.4µs ± 1%    ~     (p=0.756 n=19+16)
PushBackList-2        156µs ± 1%     158µs ± 3%  +1.03%  (p=0.024 n=19+20)
PushBackChannel-2     110µs ± 2%     110µs ± 2%    ~     (p=0.341 n=20+20)
RandomQueue-2         158µs ± 4%     161µs ± 2%  +2.34%  (p=0.000 n=20+20)
RandomList-2          279µs ± 1%     281µs ± 4%    ~     (p=0.223 n=19+19)
GrowShrinkQueue-2     111µs ± 2%     110µs ± 1%  -1.39%  (p=0.000 n=18+19)
GrowShrinkList-2      168µs ± 1%     170µs ± 5%  +1.21%  (p=0.017 n=19+20)

name               old alloc/op   new alloc/op   delta
PushFrontQueue-2     40.9kB ± 0%    40.9kB ± 0%    ~     (all equal)
PushFrontList-2      57.4kB ± 0%    57.4kB ± 0%    ~     (all equal)
PushBackQueue-2      40.9kB ± 0%    40.9kB ± 0%    ~     (all equal)
PushBackList-2       57.4kB ± 0%    57.4kB ± 0%    ~     (all equal)
PushBackChannel-2    24.7kB ± 0%    24.7kB ± 0%    ~     (all equal)
RandomQueue-2        45.7kB ± 0%    45.7kB ± 0%    ~     (all equal)
RandomList-2         90.8kB ± 0%    90.8kB ± 0%    ~     (all equal)
GrowShrinkQueue-2    57.2kB ± 0%    57.2kB ± 0%    ~     (all equal)
GrowShrinkList-2     57.4kB ± 0%    57.4kB ± 0%    ~     (all equal)

name               old allocs/op  new allocs/op  delta
PushFrontQueue-2      1.03k ± 0%     1.03k ± 0%    ~     (all equal)
PushFrontList-2       2.05k ± 0%     2.05k ± 0%    ~     (all equal)
PushBackQueue-2       1.03k ± 0%     1.03k ± 0%    ~     (all equal)
PushBackList-2        2.05k ± 0%     2.05k ± 0%    ~     (all equal)
PushBackChannel-2     1.03k ± 0%     1.03k ± 0%    ~     (all equal)
RandomQueue-2         1.63k ± 0%     1.63k ± 0%    ~     (all equal)
RandomList-2          3.24k ± 0%     3.24k ± 0%    ~     (all equal)
GrowShrinkQueue-2     1.04k ± 0%     1.04k ± 0%    ~     (all equal)
GrowShrinkList-2      2.05k ± 0%     2.05k ± 0%    ~     (all equal)
$ go version
go version go1.8.1 linux/amd64
```

The additional test cases still have *some* effect on the benchmarks, but
it's nowhere near as extreme as it was before.
I am not perfectly happy with this, but I am happy enough.
