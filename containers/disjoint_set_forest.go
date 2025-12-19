package containers

type DisjointSetForest[T comparable] struct {
	forest map[T]T
	sets   int
}

func NewDisjointSetForest[T comparable]() *DisjointSetForest[T] {
	return &DisjointSetForest[T]{
		forest: make(map[T]T),
	}
}

func (dsf *DisjointSetForest[T]) NewSet(item T) {
	if _, exists := dsf.forest[item]; !exists {
		dsf.forest[item] = item
		dsf.sets++
	}
}

func (dsf *DisjointSetForest[T]) SetCount() int {
	return dsf.sets
}

func (dsf *DisjointSetForest[T]) Find(item T) T {
	parent := dsf.forest[item]
	if parent != item {
		dsf.forest[item] = dsf.Find(parent)
	}

	return dsf.forest[item]
}

func (dsf *DisjointSetForest[T]) Union(a, b T) bool {
	aParent := dsf.Find(a)
	bParent := dsf.Find(b)

	if aParent == bParent {
		return false
	}

	dsf.forest[bParent] = aParent
	dsf.sets--

	return true
}
