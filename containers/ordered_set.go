package containers

import (
	"fmt"
	"iter"
	"slices"
)

type OrderedSet[T comparable] struct {
	set   Set[T]
	items []T
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		set: NewSet[T](),
	}
}

func (os *OrderedSet[T]) Add(item T) {
	if os.set.Has(item) {
		return
	}

	os.set.Add(item)
	os.items = append(os.items, item)
}

func (os *OrderedSet[T]) Prepend(item T) {
	os.set.Add(item)
	os.items = append([]T{item}, os.items...)
}

func (os OrderedSet[T]) Has(item T) bool {
	return os.set.Has(item)
}

func (os *OrderedSet[T]) Remove(item T) {
	if !os.set.Has(item) {
		return
	}

	os.set.Remove(item)
	index := slices.IndexFunc(os.items, func(existing T) bool {
		return existing == item
	})

	os.items = append(os.items[0:index], os.items[index+1:]...)
}

func (os OrderedSet[T]) Len() int {
	return len(os.items)
}

func (os OrderedSet[T]) At(index int) T {
	if index >= 0 {
		return os.items[index]
	}

	actualIndex := max(0, len(os.items)+index)
	return os.items[actualIndex]
}

func (os OrderedSet[T]) Iter() iter.Seq[T] {
	return slices.Values(os.items)
}

func (os OrderedSet[T]) Clone() *OrderedSet[T] {
	newItems := make([]T, len(os.items))
	copy(newItems, os.items)

	return &OrderedSet[T]{
		set:   os.set.Clone(),
		items: newItems,
	}
}

func (os OrderedSet[T]) String() string {
	return fmt.Sprint(os.items)
}
