# Queue data structure for Go

I was hacking a breadth-first search in Go and needed a queue but
all I could find in the standard library was container/list.

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
to follow even though I disagree with several of them.

## Latest Performance Comparison

The benchmarks are not very sophisticated yet but it seems that we
rather clearly beat container/list on the most common operations.

```
$ go test -bench . -benchmem
PASS
BenchmarkPushFrontQueue	20000000	       186 ns/op	      53 B/op	       0 allocs/op
BenchmarkPushFrontList	 5000000	       302 ns/op	      49 B/op	       1 allocs/op
BenchmarkPushBackQueue	20000000	       167 ns/op	      53 B/op	       0 allocs/op
BenchmarkPushBackList	 5000000	       305 ns/op	      49 B/op	       1 allocs/op
BenchmarkRandomQueue	 5000000	       422 ns/op	      26 B/op	       0 allocs/op
BenchmarkRandomList	 2000000	       797 ns/op	      78 B/op	       1 allocs/op
ok  	github.com/phf/go-queue	16.806s
```

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

In any case, presumably the Go guys know better what they want their
libraries to look like than I do, so for this queue module I simply
followed their conventions.
I would much prefer to panic in your face when you try to remove or
even just access something from an empty queue.
But since their stuff doesn't panic in similar circumstances, this
queue implementation doesn't either.
It just silently ignores the problem and hands you a nil value instead.
Now of course you have to keep checking that return value all the time
instead of being able to rely on a runtime check.
Oh well.
