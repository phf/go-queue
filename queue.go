// Copyright (c) 2013, Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

// Package queue implements a double-ended queue abstraction on
// top of a slice/array. All operations are constant time except
// for PushFront and PushBack which are amortized constant time.
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
	// length, so starting at 1 requires special cases
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
// something Go apparently likes a lot.
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
	big := make([]interface{}, q.length*2)
	j := q.front
	for i := 0; i < q.length; i++ {
		big[i] = q.rep[j]
		j = q.inc(j)
	}
	q.rep = big
	q.front = 0
	q.back = q.length
}

// TODO: leave this in or not?

func (q *Queue) String() string {
	//	result := fmt.Sprintf("(f: %d b: %d l:%d c:%d)", q.front, q.back, q.length, len(q.rep))
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
	q.lazyInit() // TODO: keep?
	if q.full() {
		q.grow()
	}
	q.front = q.dec(q.front)
	q.rep[q.front] = v
	q.length++
}

// PushBack inserts a new value v at the back of queue q.
func (q *Queue) PushBack(v interface{}) {
	q.lazyInit() // TODO: keep?
	if q.full() {
		q.grow()
	}
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
	q.rep[q.front] = nil // nice to GC?
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
	q.rep[q.back] = nil // nice to GC?
	q.length--
	return v
}
