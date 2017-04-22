// Copyright (c) 2013-2017, Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

// Package queue implements a double-ended queue abstraction on
// top of a slice/array. All operations are constant time except
// for PushFront and PushBack which are amortized constant time.
//
// We are about 15%-45% faster than container/list at the price
// of potentially wasting some memory because we grow by doubling.
// We seem to even beat Go's channels by a small margin.
package queue

import "fmt"

// Queue represents a double-ended queue.
// The zero value for Queue is an empty queue ready to use.
type Queue struct {
	// PushBack writes to rep[back] and then increments
	// back; PushFront decrements front and then writes
	// to rep[front]; gotta love those invariants.
	rep    []interface{}
	front  int
	back   int
	length int
}

// New returns an initialized empty queue.
func New() *Queue {
	return new(Queue).Init()
}

// Init initializes or clears queue q.
func (q *Queue) Init() *Queue {
	// start with a slice of length 2 even if that "wastes"
	// some memory; we do front/back arithmetic modulo the
	// length, so starting at 1 would require special cases
	q.rep = make([]interface{}, 2)
	// for some time I considered reusing the existing slice
	// if all a client does is re-initialize the queue; the
	// big problem with that is that the previous queue might
	// have been huge while the current queue doesn't grow
	// much at all; if that were to happen we'd hold on to a
	// huge chunk of memory for just a few elements and nobody
	// could do anything about it; so instead I decided to
	// just allocate a new slice and let the GC take care of
	// the previous one; seems a better tradeoff all around
	q.front, q.back, q.length = 0, 0, 0
	return q
}

// lazyInit lazily initializes a zero Queue value.
//
// I am mostly doing this because container/list does the same thing.
// Personally I think it's a little wasteful because every single
// PushFront/PushBack is going to pay the overhead of calling this.
// But that's the price for making zero values useful immediately,
// something Go folks apparently like a lot.
func (q *Queue) lazyInit() {
	if q.rep == nil {
		q.Init()
	}
}

// Len returns the number of elements of queue q.
func (q *Queue) Len() int {
	return q.length
}

// empty returns true if the queue q has no elements.
func (q *Queue) empty() bool {
	return q.length == 0
}

// full returns true if the queue q is at capacity.
func (q *Queue) full() bool {
	return q.length == len(q.rep)
}

// grow doubles the size of queue q's underlying slice/array.
func (q *Queue) grow() {
	bigger := make([]interface{}, q.length*2)
	// Kudos to Rodrigo Moraes, see https://gist.github.com/moraes/2141121
	copy(bigger, q.rep[q.front:])
	copy(bigger[q.length-q.front:], q.rep[:q.front])
	// The above replaced the "obvious" for loop and is a bit tricky.
	// First note that q.front == q.back if we're full; if that wasn't
	// true, things would be more complicated. Second recall that for
	// a slice [lo:hi] the lo bound is inclusive whereas the hi bound
	// is exclusive. If that doesn't convince you that the above works
	// maybe drawing out some pictures for a concrete example will?
	q.rep = bigger
	q.front = 0
	q.back = q.length
}

// lazyGrow grows the underlying slice/array if necessary.
func (q *Queue) lazyGrow() {
	if q.full() {
		q.grow()
	}
}

// String returns a string representation of queue q formatted
// from front to back.
func (q *Queue) String() string {
	result := ""
	result = result + "["
	j := q.front
	for i := 0; i < q.length; i++ {
		if i == q.length-1 {
			result = result + fmt.Sprintf("%v", q.rep[j])
		} else {
			result = result + fmt.Sprintf("%v, ", q.rep[j])
		}
		j = q.inc(j)
	}
	result = result + "]"
	return result
}

// inc returns the next integer position wrapping around queue q.
func (q *Queue) inc(i int) int {
	l := len(q.rep)
	return (i + 1 + l) % l
}

// dec returns the previous integer position wrapping around queue q.
func (q *Queue) dec(i int) int {
	l := len(q.rep)
	return (i - 1 + l) % l
}

// Front returns the first element of queue q or nil.
func (q *Queue) Front() interface{} {
	if q.empty() {
		return nil
	}
	return q.rep[q.front]
}

// Back returns the last element of queue q or nil.
func (q *Queue) Back() interface{} {
	if q.empty() {
		return nil
	}
	return q.rep[q.dec(q.back)]
}

// PushFront inserts a new value v at the front of queue q.
func (q *Queue) PushFront(v interface{}) {
	q.lazyInit()
	q.lazyGrow()
	q.front = q.dec(q.front)
	q.rep[q.front] = v
	q.length++
}

// PushBack inserts a new value v at the back of queue q.
func (q *Queue) PushBack(v interface{}) {
	q.lazyInit()
	q.lazyGrow()
	q.rep[q.back] = v
	q.back = q.inc(q.back)
	q.length++
}

// Both PopFront and PopBack set the newly free slot to nil
// in an attempt to be nice to the garbage collector.

// PopFront removes and returns the first element of queue q or nil.
func (q *Queue) PopFront() interface{} {
	if q.empty() {
		return nil
	}
	v := q.rep[q.front]
	q.rep[q.front] = nil
	q.front = q.inc(q.front)
	q.length--
	return v
}

// PopBack removes and returns the last element of queue q or nil.
func (q *Queue) PopBack() interface{} {
	if q.empty() {
		return nil
	}
	q.back = q.dec(q.back)
	v := q.rep[q.back]
	q.rep[q.back] = nil
	q.length--
	return v
}
