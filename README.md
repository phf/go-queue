# Queue data structure for Go

[![GoDoc](https://godoc.org/github.com/phf/go-queue/queue?status.png)](http://godoc.org/github.com/phf/go-queue/queue)
[![Go Report Card](https://goreportcard.com/badge/github.com/phf/go-queue)](https://goreportcard.com/report/github.com/phf/go-queue)

## Background

I was hacking a breadth-first search in Go and needed a queue but
all I could find in the standard library was
[container/list](https://golang.org/pkg/container/list/).

Now in principle there's nothing wrong with container/list, but I
had just admonished my students to always carefully think about
the number of memory allocations their programs make.
In other words, it felt a bit wrong for me to use a data structure
that will allocate memory for every single vertex we visit during
a breadth-first search.

So I quickly hacked a simple queue on top of a slice and finished
my project.
Now I am trying to clean up the code I wrote to give everybody else
what I really wanted the standard library to have:
A queue abstraction that doesn't allocate memory on every single
insertion.

I am trying to stick close to the conventions container/list seems
to follow even though I disagree with several of them (see below).

## Performance comparison

The benchmarks are not very sophisticated yet but we seem to beat
container/list by 15%-45% or so.
We now even beat Go's channels by about 5%, but I am still a bit
suspicious of that.
Anyway, here are the numbers:

```
$ go test -bench . -benchmem
BenchmarkPushFrontQueue-2    	   10000	    109021 ns/op	   40736 B/op	    1010 allocs/op
BenchmarkPushFrontList-2     	   10000	    154218 ns/op	   56048 B/op	    2001 allocs/op
BenchmarkPushBackQueue-2     	   10000	    107951 ns/op	   40736 B/op	    1010 allocs/op
BenchmarkPushBackList-2      	   10000	    158219 ns/op	   56048 B/op	    2001 allocs/op
BenchmarkPushBackChannel-2   	   10000	    113839 ns/op	   24480 B/op	    1002 allocs/op
BenchmarkRandomQueue-2       	    2000	    600828 ns/op	   45530 B/op	    1610 allocs/op
BenchmarkRandomList-2        	    2000	    692993 ns/op	   89667 B/op	    3201 allocs/op
PASS
ok  	github.com/phf/go-queue/queue	9.241s
```

### Go's channels as queues

Go's channels *used* to beat our queue implementation by about 22%
for `PushBack`.
(In fact I used to call them "*ridiculously* fast" before and
recommended their use in situations where nothing but performance
matters.)
That seemed sensible considering that channels are built into the
language and offer a lot less functionality:
We have to size them correctly if we want to use them as a simple
queue in an otherwise non-concurrent setting, they are not
double-ended, and they don't support "peeking" at the next element
without removing it.
Apparently replacing the "manual" loop when a queue has to grow with
[copy](https://golang.org/ref/spec#Appending_and_copying_slices) has
paid off.
(That or I am benchmarking this incorrectly.)

## What I don't like about Go's conventions

I guess my biggest gripe with Go's container/list is that it tries
very hard to *never* **ever** panic.
I don't understand this, and in fact I think it's rather dangerous.

Take a plain old array for example.
When you index outside of its domain, you get a panic even in Go.
As you should!
This kind of runtime check helps you catch your indexing errors and
it also enforces the abstraction provided by the array.

But then Go already messes things up with the builtin map type.
Instead of getting a panic when you try to access a key that's not
in the map, you get a zero value.
And if you *really* want to know whether a key is there or not you
have to go through some extra stunts.

Apparently they just kept going from there with the libraries.
In the case of container/list for example, if you try to remove
an element that's *not* *actually* *from* *that* *list*, nothing
happens.
Instead of immediately getting into your face with a panic and
helping you fix your code, you'll just keep wondering why the
Remove() operation you wrote down didn't work.
Indeed you'll probably end up looking for the bug in all the wrong
places before it finally dawns on you that maybe you removed from
the wrong list.

In any case, presumably the Go folks know better what they want their
libraries to look like than I do, so for this queue module I simply
followed their conventions.
I would much prefer to panic in your face when you try to remove or
even just access something from an empty queue.
But since their stuff doesn't panic in similar circumstances, this
queue implementation doesn't either.

## Kudos

- [Rodrigo Moraes](https://github.com/moraes) for posting
  [this gist](https://gist.github.com/moraes/2141121) which reminded
  me of Go's [copy](https://golang.org/ref/spec#Appending_and_copying_slices)
  builtin and a similar trick I had previously used in Java.
