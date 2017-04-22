# Queue data structure for Go

The main repository is [here](https://github.com/phf/go-queue).
This is just a rant that really had no place in the "official"
`README` file.

## What I don't like about Go's conventions

I guess my biggest gripe with Go's
[container/list](https://golang.org/pkg/container/list/) is that it
tries *very* hard to *never* **ever** panic.
I don't understand this, and in fact I think it's rather dangerous.

Take a plain old array for example.
When you index outside of its domain, you get a panic even in Go.
As you should!
This kind of runtime check helps you catch your indexing errors and
it also enforces the abstraction provided by the array.

But then Go already "messes things up" with the builtin map type.
Instead of getting a panic when you try to access a key that's not
in the map, you get a zero value.
And if you *really* want to know whether a key is there or not you
have to go through some extra stunts.

Apparently they just kept going from there with the libraries.
In the case of [container/list](https://golang.org/pkg/container/list/)
for example, if you try to remove an element that's *not* *actually*
*from* **that** *list*, nothing happens.
Instead of immediately getting into your face with a panic and
helping you fix your code, you'll just keep wondering why the
`Remove` operation you wrote down didn't work.
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
