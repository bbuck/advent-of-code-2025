package containers

import (
	"iter"
	"maps"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	m := make(map[T]struct{})

	return Set[T](m)
}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Has(item T) bool {
	_, exists := s[item]

	return exists
}

func (s Set[T]) Iter() iter.Seq[T] {
	return maps.Keys(s)
}
