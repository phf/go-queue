# Queue data structure for Go

[![GoDoc](https://godoc.org/github.com/phf/go-queue/queue?status.png)](http://godoc.org/github.com/phf/go-queue/queue)
[![Go Report Card](https://goreportcard.com/badge/github.com/phf/go-queue)](https://goreportcard.com/report/github.com/phf/go-queue)

A double-ended queue (aka "deque") built on top of a slice.
All operations are (amortized) constant time.
Benchmarks compare favorably to
[container/list](https://golang.org/pkg/container/list/) as
well as to Go's channels.

I tried to stick to the conventions established by
[container/list](https://golang.org/pkg/container/list/)
even though I disagree with them (see
[`RANT.md`](https://github.com/phf/go-queue/blob/master/RANT.md)
for details).
In other words, this data structure is ready for the standard
library (hah!).

## Background

In 2013 I was hacking a breadth-first search in Go and needed a
queue, but all I could find in the standard library was
[container/list](https://golang.org/pkg/container/list/).

Now in *principle* there's nothing wrong with
[container/list](https://golang.org/pkg/container/list/), but I
had just admonished my students to *always* think carefully about
the number of memory allocations their programs make.
In other words, it felt wrong for me to use a data structure that
allocates memory for *every* single vertex we visit during a
breadth-first search.

After I got done with my project, I decided to clean up the queue
code a little and to push it here to give everybody else what I
really wanted to find in the standard library:
A queue abstraction that doesn't allocate memory on every single
insertion.

## Performance

**Please read
[`BENCH.md`](https://github.com/phf/go-queue/blob/master/BENCH.md)
for some perspective.
The numbers below are most likely "contaminated" in a way that makes
our queues appear *worse* than they are.**

Here are the numbers for my (ancient) home machine:

```
$ go test -bench=. -benchmem -count=10 >bench.txt
$ benchstat bench.txt
name               time/op
PushFrontQueue-2   97.7µs ± 1%
PushFrontList-2     163µs ± 1%
PushBackQueue-2    98.0µs ± 1%
PushBackList-2      165µs ± 3%
PushBackChannel-2   145µs ± 1%
RandomQueue-2       172µs ± 1%
RandomList-2        292µs ± 1%
GrowShrinkQueue-2   121µs ± 1%
GrowShrinkList-2    174µs ± 1%

name               alloc/op
PushFrontQueue-2   40.9kB ± 0%
PushFrontList-2    57.4kB ± 0%
PushBackQueue-2    40.9kB ± 0%
PushBackList-2     57.4kB ± 0%
PushBackChannel-2  24.7kB ± 0%
RandomQueue-2      45.7kB ± 0%
RandomList-2       90.8kB ± 0%
GrowShrinkQueue-2  57.2kB ± 0%
GrowShrinkList-2   57.4kB ± 0%

name               allocs/op
PushFrontQueue-2    1.03k ± 0%
PushFrontList-2     2.05k ± 0%
PushBackQueue-2     1.03k ± 0%
PushBackList-2      2.05k ± 0%
PushBackChannel-2   1.03k ± 0%
RandomQueue-2       1.63k ± 0%
RandomList-2        3.24k ± 0%
GrowShrinkQueue-2   1.04k ± 0%
GrowShrinkList-2    2.05k ± 0%
$ go version
go version go1.7.5 linux/amd64
$ cat /proc/cpuinfo | grep "model name" | uniq
model name	: AMD Athlon(tm) 64 X2 Dual Core Processor 6000+
```

That's a [speedup](https://en.wikipedia.org/wiki/Speedup) of
1.45-1.70
over [container/list](https://golang.org/pkg/container/list/) and a speedup of
1.48
over Go's channels.
We also consistently allocate less memory in fewer allocations than
[container/list](https://golang.org/pkg/container/list/).
(Note that the number of allocations seems off: since we grow by *doubling*
we should only allocate memory *O(log n)* times.)

The same benchmarks on one of our department's servers:

```
$ go test -bench=. -benchmem -count=10 >bench.txt
$ benchstat bench.txt
name               time/op
PushFrontQueue-8   88.8µs ± 3%
PushFrontList-8     156µs ± 5%
PushBackQueue-8    88.3µs ± 1%
PushBackList-8      159µs ± 2%
PushBackChannel-8   132µs ± 2%
RandomQueue-8       156µs ± 7%
RandomList-8        279µs ±10%
GrowShrinkQueue-8   117µs ± 0%
GrowShrinkList-8    164µs ± 4%

name               alloc/op
PushFrontQueue-8   40.9kB ± 0%
PushFrontList-8    57.4kB ± 0%
PushBackQueue-8    40.9kB ± 0%
PushBackList-8     57.4kB ± 0%
PushBackChannel-8  24.7kB ± 0%
RandomQueue-8      45.7kB ± 0%
RandomList-8       90.8kB ± 0%
GrowShrinkQueue-8  57.2kB ± 0%
GrowShrinkList-8   57.4kB ± 0%

name               allocs/op
PushFrontQueue-8    1.03k ± 0%
PushFrontList-8     2.05k ± 0%
PushBackQueue-8     1.03k ± 0%
PushBackList-8      2.05k ± 0%
PushBackChannel-8   1.03k ± 0%
RandomQueue-8       1.63k ± 0%
RandomList-8        3.24k ± 0%
GrowShrinkQueue-8   1.04k ± 0%
GrowShrinkList-8    2.05k ± 0%
$ go version
go version go1.7.5 linux/amd64
$ cat /proc/cpuinfo | grep "model name" |uniq
model name	: Intel(R) Xeon(R) CPU           E5440  @ 2.83GHz
```

That's a [speedup](https://en.wikipedia.org/wiki/Speedup) of
1.76-1.80
over [container/list](https://golang.org/pkg/container/list/) and a speedup of
1.49
over Go's channels.

The same benchmarks on a *different* department server:

```
$ go test -bench=. -benchmem -count=10 >bench.txt
$ benchstat bench.txt
name                time/op
PushFrontQueue-24   89.1µs ± 8%
PushFrontList-24     176µs ± 8%
PushBackQueue-24    86.8µs ± 5%
PushBackList-24      178µs ± 6%
PushBackChannel-24   151µs ±12%
RandomQueue-24       180µs ±24%
RandomList-24        334µs ± 7%
GrowShrinkQueue-24   117µs ± 3%
GrowShrinkList-24    187µs ± 6%

name                alloc/op
PushFrontQueue-24   40.9kB ± 0%
PushFrontList-24    57.4kB ± 0%
PushBackQueue-24    40.9kB ± 0%
PushBackList-24     57.4kB ± 0%
PushBackChannel-24  24.7kB ± 0%
RandomQueue-24      45.7kB ± 0%
RandomList-24       90.8kB ± 0%
GrowShrinkQueue-24  57.2kB ± 0%
GrowShrinkList-24   57.4kB ± 0%

name                allocs/op
PushFrontQueue-24    1.03k ± 0%
PushFrontList-24     2.05k ± 0%
PushBackQueue-24     1.03k ± 0%
PushBackList-24      2.05k ± 0%
PushBackChannel-24   1.03k ± 0%
RandomQueue-24       1.63k ± 0%
RandomList-24        3.24k ± 0%
GrowShrinkQueue-24   1.04k ± 0%
GrowShrinkList-24    2.05k ± 0%
$ go version
go version go1.7.4 linux/amd64
$ cat /proc/cpuinfo | grep "model name" |uniq
model name	: Intel(R) Xeon(R) CPU E5-2420 0 @ 1.90GHz
```

That's a [speedup](https://en.wikipedia.org/wiki/Speedup) of
1.86-**2.05**
over [container/list](https://golang.org/pkg/container/list/) and a speedup of
1.74
over Go's channels.

The same benchmarks on an old
[Raspberry Pi Model B Rev 1](https://en.wikipedia.org/wiki/Raspberry_Pi):

```
$ benchstat bench.txt
name             time/op
PushFrontQueue    788µs ±24%
PushFrontList    2.74ms ±14%
PushBackQueue    1.11ms ± 3%
PushBackList     2.73ms ±14%
PushBackChannel  1.25ms ± 3%
RandomQueue      1.50ms ± 1%
RandomList       4.92ms ± 6%
GrowShrinkQueue  1.26ms ± 0%
GrowShrinkList   2.88ms ± 2%

name             alloc/op
PushFrontQueue   16.5kB ± 0%
PushFrontList    33.9kB ± 0%
PushBackQueue    16.5kB ± 0%
PushBackList     33.9kB ± 0%
PushBackChannel  8.45kB ± 0%
RandomQueue      16.5kB ± 0%
RandomList       53.4kB ± 0%
GrowShrinkQueue  24.6kB ± 0%
GrowShrinkList   33.9kB ± 0%

name             allocs/op
PushFrontQueue     12.0 ± 0%
PushFrontList     1.03k ± 0%
PushBackQueue      12.0 ± 0%
PushBackList      1.03k ± 0%
PushBackChannel    1.00 ± 0%
RandomQueue        12.0 ± 0%
RandomList        1.63k ± 0%
GrowShrinkQueue    20.0 ± 0%
GrowShrinkList    1.03k ± 0%
$ go version
go version go1.3.3 linux/arm
$ cat /proc/cpuinfo |grep "model name"
model name	: ARMv6-compatible processor rev 7 (v6l)
```

That's a [speedup](https://en.wikipedia.org/wiki/Speedup) of
**2.46-3.48**
over [container/list](https://golang.org/pkg/container/list/)
but only a speedup of
1.13
over Go's channels.
(Note that I had to manually repeat the benchmarks and then run `benchtest`
elsewhere since those features/tools are not available for Go 1.3;
however, the number of allocations seems to be correct here for the first
time, maybe there's some breakage in the more recent benchmarking
framework?)

### Go's channels as queues

Go's channels *used* to beat our queue implementation by about 22%
for `PushBack`.
That seemed sensible considering that channels are built into the
language and offer a lot less functionality:
We have to size them correctly if we want to use them as a simple
queue in an otherwise non-concurrent setting, they are not
double-ended, and they don't support "peeking" at the next element
without removing it.

That all changed with
[two](https://github.com/phf/go-queue/commit/5652cbe39198516d853918fe64a4e70948b42f1a)
[commits](https://github.com/phf/go-queue/commit/aa6086b89f98eb5cfd8df918e57612271ae1c137)
in which I replaced the "manual" loop when a queue has to grow with
[copy](https://golang.org/ref/spec#Appending_and_copying_slices)
and the `%` operations to wrap indices around the slice with
equivalent `&` operations.
(The code was originally written without these "hacks" because I wanted to
show it to my "innocent" Java students.)
Those two changes *really* paid off.

(I used to call channels "*ridiculously* fast" before and recommended their
use in situations where nothing but performance matters.
Alas that may no longer be good advice.
Either that, or I am just benchmarking incorrectly.)

## Kudos

Hacking queue data structures in Go seems to be a popular way to spend
an evening. Kudos to...

- [Rodrigo Moraes](https://github.com/moraes) for posting
  [this gist](https://gist.github.com/moraes/2141121) which reminded
  me of Go's [copy](https://golang.org/ref/spec#Appending_and_copying_slices)
  builtin and a similar trick I had previously used in Java.
- [Evan Huus](https://github.com/eapache) for sharing
  [his queue](https://github.com/eapache/queue) which reminded me of
  the old "replace % by &" trick I had used many times before.
- [Dariusz Górecki](https://github.com/canni) for his
  [commit](https://github.com/eapache/queue/commit/334cc1b02398be651373851653017e6cbf588f9e)
  to [Evan](https://github.com/eapache)'s queue that simplified
  [Rodrigo](https://github.com/moraes)'s snippet and hence mine.

If you find something in my code that helps you improve yours, feel
free to run with it!
