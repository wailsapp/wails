// Code from https://github.com/erikdubbelboer/ringqueue
/*
The MIT License (MIT)

Copyright (c) 2015 Erik Dubbelboer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package assetserver

type ringqueue[T any] struct {
	nodes []T
	head  int
	tail  int
	cnt   int

	minSize int
}

func newRingqueue[T any](minSize uint) *ringqueue[T] {
	if minSize < 2 {
		minSize = 2
	}
	return &ringqueue[T]{
		nodes:   make([]T, minSize),
		minSize: int(minSize),
	}
}

func (q *ringqueue[T]) resize(n int) {
	nodes := make([]T, n)
	if q.head < q.tail {
		copy(nodes, q.nodes[q.head:q.tail])
	} else {
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.tail])
	}

	q.tail = q.cnt % n
	q.head = 0
	q.nodes = nodes
}

func (q *ringqueue[T]) Add(i T) {
	if q.cnt == len(q.nodes) {
		// Also tested a grow rate of 1.5, see: http://stackoverflow.com/questions/2269063/buffer-growth-strategy
		// In Go this resulted in a higher memory usage.
		q.resize(q.cnt * 2)
	}
	q.nodes[q.tail] = i
	q.tail = (q.tail + 1) % len(q.nodes)
	q.cnt++
}

func (q *ringqueue[T]) Peek() (T, bool) {
	if q.cnt == 0 {
		var none T
		return none, false
	}
	return q.nodes[q.head], true
}

func (q *ringqueue[T]) Remove() (T, bool) {
	if q.cnt == 0 {
		var none T
		return none, false
	}
	i := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.cnt--

	if n := len(q.nodes) / 2; n > q.minSize && q.cnt <= n {
		q.resize(n)
	}

	return i, true
}

func (q *ringqueue[T]) Cap() int {
	return cap(q.nodes)
}

func (q *ringqueue[T]) Len() int {
	return q.cnt
}
