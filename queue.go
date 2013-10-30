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
	rep []interface{}
	front, back, length int
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
	q.rep = make([]interface{}, 2) // TODO: keep old slice/array? but what if it's huge?
	q.front, q.back, q.length = 0, 0, 0
	return q
}

// TODO: good idea? list.go does it but slows down every insertion...
// I guess that's the price for allowing zero values to be useful?
func (q *Queue) lazyInit() {
	if q.rep == nil {
		q.Init()
	}
}

// Len returns the number of elements of queue q.
func (q *Queue) Len() int {
	return q.length
}

func (q *Queue) empty() bool {
	return q.length == 0
}

func (q *Queue) full() bool {
	return q.length == len(q.rep)
}

func (q *Queue) grow() {
	big := make([]interface{}, q.length*2)
	j := q.front
	for i := 0; i < q.length; i++ {
		big[i] = q.rep[j]
		q.inc(&j)
	}
	q.rep = big
	q.front = 0
	q.back = q.length
}

// TODO: leave this in or not?
func (q *Queue) String() string {
	result := fmt.Sprintf("(f: %d b: %d l:%d c:%d)", q.front, q.back, q.length, len(q.rep))
	result = result + "["
	j := q.front
	for i := 0; i < q.length; i++ {
		result = result + fmt.Sprintf("[%v]", q.rep[j])
		q.inc(&j)
	}
	result = result + "]"
	return result
}

// TODO: convert these two back to proper functions? see ugliness in Back() below

func (q *Queue) inc(i *int) {
	l := len(q.rep)
	*i = (*i+1+l) % l
}

func (q *Queue) dec(i *int) {
	l := len(q.rep)
	*i = (*i-1+l) % l
}

// TODO: I dislike the Go philosophy of avoiding panics at all
// costs; Front/Back/Pop from an empty Queue SHOULD panic! at
// least in my mind...

// Front returns the first element of queue q or nil. 
func (q *Queue) Front() interface{} {
	if q.empty() { return nil }
	return q.rep[q.front]
}

// Back returns the last element of queue q or nil. 
func (q *Queue) Back() interface{} {
	if q.empty() { return nil }
	b := q.back
	q.dec(&b)
	return q.rep[b]
}

// PushFront inserts a new value v at the front of queue q.
func (q *Queue) PushFront(v interface{}) {
	q.lazyInit() // TODO: keep?
	if q.full() { q.grow() }
	q.dec(&q.front)
	q.rep[q.front] = v
	q.length++
}

// PushBack inserts a new value v at the back of queue q.
func (q *Queue) PushBack(v interface{}) {
	q.lazyInit() // TODO: keep?
	if q.full() { q.grow() }
	q.rep[q.back] = v
	q.inc(&q.back)
	q.length++
}

// PopFront removes and returns the first element of queue q or nil.
func (q *Queue) PopFront() interface{} {
	if q.empty() { return nil }
	v := q.rep[q.front]
	q.rep[q.front] = nil // nice to GC?
	q.inc(&q.front)
	q.length--
	return v
}

// PopBack removes and returns the last element of queue q or nil.
func (q *Queue) PopBack() interface{} {
	if q.empty() { return nil }
	q.dec(&q.back)
	v := q.rep[q.back]
	q.rep[q.back] = nil // nice to GC?
	q.length--
	return v
}
