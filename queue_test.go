// Copyright (c) 2013, Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package queue

// TODO: need a lot more tests, and maybe a better way of
// modularizing them; also benchmarks comparing this to
// Go's container.list

import "testing"
import "container/list"
import "math/rand"

func ensureEmpty(t *testing.T, q *Queue) {
	if l := q.Len(); l != 0 {
		t.Errorf("q.Len() = %d, want %d", l, 0)
	}
	if e := q.Front(); e != nil {
		t.Errorf("q.Front() = %v, want %v", e, nil)
	}
	if e := q.Back(); e != nil {
		t.Errorf("q.Back() = %v, want %v", e, nil)
	}
}

func TestNew(t *testing.T) {
	q := New()
	ensureEmpty(t, q)
}

func ensureSingleton(t *testing.T, q *Queue) {
	if l := q.Len(); l != 1 {
		t.Errorf("q.Len() = %d, want %d", l, 1)
	}
	if e := q.Front(); e != 42 {
		t.Errorf("q.Front() = %v, want %v", e, 42)
	}
	if e := q.Back(); e != 42 {
		t.Errorf("q.Back() = %v, want %v", e, 42)
	}
}

func TestSingleton(t *testing.T) {
	q := New()
	ensureEmpty(t, q)
	q.PushFront(42)
	ensureSingleton(t, q)
	q.PopFront()
	ensureEmpty(t, q)
	q.PushBack(42)
	ensureSingleton(t, q)
	q.PopBack()
	ensureEmpty(t, q)
	q.PushFront(42)
	ensureSingleton(t, q)
	q.PopBack()
	ensureEmpty(t, q)
	q.PushBack(42)
	ensureSingleton(t, q)
	q.PopFront()
	ensureEmpty(t, q)
}

func TestDuos(t *testing.T) {
	q := New()
	ensureEmpty(t, q)
	q.PushFront(42)
	ensureSingleton(t, q)
	q.PushBack(43)
	if l := q.Len(); l != 2 {
		t.Errorf("q.Len() = %d, want %d", l, 2)
	}
	if e := q.Front(); e != 42 {
		t.Errorf("q.Front() = %v, want %v", e, 42)
	}
	if e := q.Back(); e != 43 {
		t.Errorf("q.Back() = %v, want %v", e, 43)
	}
}

func ensureLength(t *testing.T, q *Queue, len int) {
	if l := q.Len(); l != len {
		t.Errorf("q.Len() = %d, want %d", l, len)
	}
}

func TestZeroValue(t *testing.T) {
	var q Queue
	q.PushFront(1)
	ensureLength(t, &q, 1)
	q.PushFront(2)
	ensureLength(t, &q, 2)
	q.PushFront(3)
	ensureLength(t, &q, 3)
	q.PushFront(4)
	ensureLength(t, &q, 4)
	q.PushFront(5)
	ensureLength(t, &q, 5)
}

func BenchmarkPushFrontQueue(b *testing.B) {
	var q Queue
	for i := 0; i < b.N; i++ {
		q.PushFront(i)
	}
}
func BenchmarkPushFrontList(b *testing.B) {
	var q list.List
	for i := 0; i < b.N; i++ {
		q.PushFront(i)
	}
}

func BenchmarkPushBackQueue(b *testing.B) {
	var q Queue
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
}
func BenchmarkPushBackList(b *testing.B) {
	var q list.List
	for i := 0; i < b.N; i++ {
		q.PushBack(i)
	}
}

func BenchmarkRandomQueue(b *testing.B) {
	var q Queue
	rand.Seed(64738)
	for i := 0; i < b.N; i++ {
		if rand.Float32() < 0.8 {
			q.PushBack(i)
		}
		if rand.Float32() < 0.8 {
			q.PushFront(i)
		}
		if rand.Float32() < 0.5 {
			q.PopFront()
		}
		if rand.Float32() < 0.5 {
			q.PopBack()
		}
	}
}
func BenchmarkRandomList(b *testing.B) {
	var q list.List
	rand.Seed(64738)
	for i := 0; i < b.N; i++ {
		if rand.Float32() < 0.8 {
			q.PushBack(i)
		}
		if rand.Float32() < 0.8 {
			q.PushFront(i)
		}
		if rand.Float32() < 0.5 {
			if e := q.Front(); e != nil {
				q.Remove(e)
			}
		}
		if rand.Float32() < 0.5 {
			if e := q.Back(); e != nil {
				q.Remove(e)
			}
		}
	}
}
