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
