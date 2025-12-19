package containers

import (
	"container/heap"
	"iter"
	"slices"
)

type Heap[T any] struct {
	items []T
	less  func(a, b T) bool
}

func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		less: less,
	}
}

func (h *Heap[T]) Add(item T) {
	heap.Push(h, item)
}

func (h *Heap[T]) Remove() (T, bool) {
	if len(h.items) == 0 {
		var zero T

		return zero, false
	}

	return heap.Pop(h).(T), true
}

func (h *Heap[T]) Peek() (T, bool) {
	if len(h.items) == 0 {
		var zero T

		return zero, false
	}

	return h.items[0], true
}

func (h *Heap[T]) Iter() iter.Seq[T] {
	return slices.Values(h.items)
}

// containers/heap.Interface implementation

func (h *Heap[T]) Len() int {
	return len(h.items)
}

func (h *Heap[T]) Less(i, j int) bool {
	return h.less(h.items[i], h.items[j])
}

func (h *Heap[T]) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *Heap[T]) Push(item any) {
	h.items = append(h.items, item.(T))
}

func (h *Heap[T]) Pop() any {
	length := len(h.items)
	out := h.items[length-1]
	h.items = h.items[0 : length-1]

	return out
}
