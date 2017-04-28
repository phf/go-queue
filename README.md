# Queue data structure for Go

[![GoDoc](https://godoc.org/github.com/phf/go-queue/queue?status.png)](http://godoc.org/github.com/phf/go-queue/queue)
[![Go Report Card](https://goreportcard.com/badge/github.com/phf/go-queue)](https://goreportcard.com/report/github.com/phf/go-queue)

A double-ended queue (aka "deque") built on top of a slice.
All operations except pushes are constant-time; pushes are
*amortized* constant-time.
Benchmarks compare favorably to
[container/list](https://golang.org/pkg/container/list/) as
well as to Go's channels.

I tried to stick close to the conventions
[container/list](https://golang.org/pkg/container/list/) seems to
follow even though I disagree with several of them (see
[`RANT.md`](https://github.com/phf/go-queue/blob/master/RANT.md)).
In other words, it's ready for the standard library (hah!).

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

The benchmarks are not very sophisticated but we seem to be *almost*
twice as fast as [container/list](https://golang.org/pkg/container/list/)
([speedup](https://en.wikipedia.org/wiki/Speedup) of 1.85-1.93).
We're also a bit faster than Go's channels (speedup of 1.38).
Anyway, here are some numbers from my old home machine:

```
$ go test -bench . -benchmem
BenchmarkPushFrontQueue-2    	   20000	     85886 ns/op	   40944 B/op	    1035 allocs/op
BenchmarkPushFrontList-2     	   10000	    158998 ns/op	   57392 B/op	    2049 allocs/op
BenchmarkPushBackQueue-2     	   20000	     85189 ns/op	   40944 B/op	    1035 allocs/op
BenchmarkPushBackList-2      	   10000	    160718 ns/op	   57392 B/op	    2049 allocs/op
BenchmarkPushBackChannel-2   	   10000	    117610 ns/op	   24672 B/op	    1026 allocs/op
BenchmarkRandomQueue-2       	   10000	    144867 ns/op	   45720 B/op	    1632 allocs/op
BenchmarkRandomList-2        	    5000	    278965 ns/op	   90824 B/op	    3243 allocs/op
PASS
ok  	github.com/phf/go-queue/queue	12.472s
$ go version
go version go1.7.5 linux/amd64
$ uname -p
AMD Athlon(tm) 64 X2 Dual Core Processor 6000+
$ date
Sat Apr 22 11:26:40 EDT 2017
```

(The number of allocations seems off, since we grow by doubling we should
only allocate memory O(log n) times.)
The same benchmarks on a more recent laptop:

```
$ go test -bench=. -benchmem
PASS
BenchmarkPushFrontQueue-4 	   10000	    107377 ns/op	   40944 B/op	    1035 allocs/op
BenchmarkPushFrontList-4  	   10000	    205141 ns/op	   57392 B/op	    2049 allocs/op
BenchmarkPushBackQueue-4  	   10000	    107339 ns/op	   40944 B/op	    1035 allocs/op
BenchmarkPushBackList-4   	   10000	    204100 ns/op	   57392 B/op	    2049 allocs/op
BenchmarkPushBackChannel-4	   10000	    174319 ns/op	   24672 B/op	    1026 allocs/op
BenchmarkRandomQueue-4    	   10000	    190498 ns/op	   45720 B/op	    1632 allocs/op
BenchmarkRandomList-4     	    5000	    364802 ns/op	   90825 B/op	    3243 allocs/op
ok  	github.com/phf/go-queue/queue	11.881s
$ go version
go version go1.6.2 linux/amd64
$ cat /proc/cpuinfo | grep "model name" | uniq
model name	: AMD A10-4600M APU with Radeon(tm) HD Graphics
$ date
Fri Apr 28 17:20:57 EDT 2017
```

So that's a [speedup](https://en.wikipedia.org/wiki/Speedup) of 1.90 over
[container/list](https://golang.org/pkg/container/list/) and of 1.62 over
Go's channels.
The same benchmarks on an old
[Raspberry Pi Model B Rev 1](https://en.wikipedia.org/wiki/Raspberry_Pi):

```
$ go test -bench . -benchmem
PASS
BenchmarkPushFrontQueue     2000            788316 ns/op           16469 B/op         12 allocs/op
BenchmarkPushFrontList      1000           2629835 ns/op           33904 B/op       1028 allocs/op
BenchmarkPushBackQueue      2000            776663 ns/op           16469 B/op         12 allocs/op
BenchmarkPushBackList       1000           2817162 ns/op           33877 B/op       1028 allocs/op
BenchmarkPushBackChannel    2000           1229474 ns/op            8454 B/op          1 allocs/op
BenchmarkRandomQueue        2000           1325947 ns/op           16469 B/op         12 allocs/op
BenchmarkRandomList          500           4929491 ns/op           53437 B/op       1627 allocs/op
ok      github.com/phf/go-queue/queue   17.798s
$ go version
go version go1.3.3 linux/arm
$ cat /proc/cpuinfo | grep "model name"
model name      : ARMv6-compatible processor rev 7 (v6l)
$ date
Sat Apr 22 18:04:16 UTC 2017
```

So that's a [speedup](https://en.wikipedia.org/wiki/Speedup) of
**3.34**-**3.72** over
[container/list](https://golang.org/pkg/container/list/) and of 1.58 over
Go's channels.
(Also the number of allocations seems to be correct here for some
reason?)

### Go's channels as queues

Go's channels *used* to beat our queue implementation by about 22%
for `PushBack`.
That seemed sensible considering that channels are built into the
language and offer a lot less functionality:
We have to size them correctly if we want to use them as a simple
queue in an otherwise non-concurrent setting, they are not
double-ended, and they don't support "peeking" at the next element
without removing it.

It all changed with
[two](https://github.com/phf/go-queue/commit/5652cbe39198516d853918fe64a4e70948b42f1a)
[commits](https://github.com/phf/go-queue/commit/aa6086b89f98eb5cfd8df918e57612271ae1c137)
that replaced the "manual" loop when a queue has to grow with
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
