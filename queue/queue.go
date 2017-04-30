// Copyright (c) 2013-2017, Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

// Package queue implements a double-ended queue (aka "deque") on top
// of a slice. All operations are (amortized) constant time.
// Benchmarks compare favorably to container/list as well as to Go's
// channels.
package queue

import (
	"bytes"
	"fmt"
)

// Queue represents a double-ended queue.
// The zero value for Queue is an empty queue ready to use.
type Queue struct {
	// PushBack writes to rep[back] and then increments
	// back; PushFront decrements front and then writes
	// to rep[front]; len(rep) must be a power of two;
	// gotta love those invariants.
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
	q.rep = make([]interface{}, 1)
	// I considered reusing the existing slice if all a client does
	// is re-initialize the queue. The problem is that the current
	// queue might be huge, but the next one might not grow much. So
	// we'd hold on to a huge chunk of memory for just a few elements
	// and nobody can do anything. Making a new slice and letting the
	// GC take care of the old one seems like a better tradeoff.
	q.front, q.back, q.length = 0, 0, 0
	return q
}

// lazyInit lazily initializes a zero Queue value.
//
// I am mostly doing this because container/list does the same thing.
// Personally I think it's a little wasteful because every single
// PushFront/PushBack is going to pay the overhead of calling this.
// But that's the price for making zero values useful immediately.
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

// sparse returns true if the queue q has excess capacity.
func (q *Queue) sparse() bool {
	return 1 < q.length && q.length < len(q.rep)/4
}

// grow doubles the size of queue q's underlying slice.
func (q *Queue) grow() {
	bigger := make([]interface{}, len(q.rep)*2)
	// Kudos to Rodrigo Moraes, see https://gist.github.com/moraes/2141121
	// Kudos to Dariusz Górecki, see https://github.com/eapache/queue/commit/334cc1b02398be651373851653017e6cbf588f9e
	n := copy(bigger, q.rep[q.front:])
	copy(bigger[n:], q.rep[:q.back])
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

// lazyGrow grows the underlying slice if necessary.
func (q *Queue) lazyGrow() {
	if q.full() {
		q.grow()
	}
}

// shrink halves the size of queue q's underlying slice.
func (q *Queue) shrink() {
	smaller := make([]interface{}, len(q.rep)/2)
	if q.front < q.back {
		copy(smaller, q.rep[q.front:q.back])
	} else {
		n := copy(smaller, q.rep[q.front:])
		copy(smaller[n:], q.rep[:q.back])
	}
	q.rep = smaller
	q.front = 0
	q.back = q.length
}

// lazyShrink shrinks the underlying slice if advisable.
func (q *Queue) lazyShrink() {
	if q.sparse() {
		q.shrink()
	}
}

// String returns a string representation of queue q formatted
// from front to back.
func (q *Queue) String() string {
	var result bytes.Buffer
	result.WriteString("[")
	j := q.front
	for i := 0; i < q.length; i++ {
		result.WriteString(fmt.Sprintf("%v", q.rep[j]))
		if i < q.length-1 {
			result.WriteRune(' ')
		}
		j = q.inc(j)
	}
	result.WriteString("]")
	return result.String()
}

// inc returns the next integer position wrapping around queue q.
func (q *Queue) inc(i int) int {
	return (i + 1) & (len(q.rep) - 1) // requires l = 2^n
}

// dec returns the previous integer position wrapping around queue q.
func (q *Queue) dec(i int) int {
	return (i - 1) & (len(q.rep) - 1) // requires l = 2^n
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

// PopFront removes and returns the first element of queue q or nil.
func (q *Queue) PopFront() interface{} {
	if q.empty() {
		return nil
	}
	v := q.rep[q.front]
	q.rep[q.front] = nil // be nice to GC
	q.front = q.inc(q.front)
	q.length--
	q.lazyShrink()
	return v
}

// PopBack removes and returns the last element of queue q or nil.
func (q *Queue) PopBack() interface{} {
	if q.empty() {
		return nil
	}
	q.back = q.dec(q.back)
	v := q.rep[q.back]
	q.rep[q.back] = nil // be nice to GC
	q.length--
	q.lazyShrink()
	return v
}
