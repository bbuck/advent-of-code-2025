package grid

import "iter"

// Grid handles building and access values in a matrix.
type Grid[T any] struct {
	matrix [][]T

	rowLen    int
	columnLen int
}

// NewGrid creates a enw grid with the specified number of rows and columns.
func NewGrid[T any](rowLen, columnLen int) *Grid[T] {
	var matrix [][]T
	for range rowLen {
		matrix = append(matrix, make([]T, columnLen))
	}

	return &Grid[T]{
		matrix:    matrix,
		rowLen:    rowLen,
		columnLen: columnLen,
	}
}

// RowLen returns the number of rows in the Grid.
func (g Grid[T]) RowLen() int {
	return g.rowLen
}

// ColumnLen returns the number of columns in the Grid.
func (g Grid[T]) ColumnLen() int {
	return g.columnLen
}

// At returns the value at the given location within the grid.
func (g Grid[T]) At(l Location) (T, bool) {
	if !g.ValidLocation(l) {
		var zero T

		return zero, false
	}

	return g.matrix[l.Row][l.Column], true
}

// SetAt will indiscriminately update the value at the specified location.
func (g Grid[T]) SetAt(l Location, value T) bool {
	if !g.ValidLocation(l) {
		return false
	}

	g.matrix[l.Row][l.Column] = value

	return true
}

// UpdateAt will run the given update function on the current value of the cell
// and assign the result to the cell at the given location.
func (g Grid[T]) UpdateAt(l Location, update func(T) T) bool {
	if !g.ValidLocation(l) {
		return false
	}

	g.matrix[l.Row][l.Column] = update(g.matrix[l.Row][l.Column])

	return true
}

// Iter returns a row, column ordered iterator over the grid.
func (g Grid[T]) Iter() iter.Seq2[Location, T] {
	return func(yield func(Location, T) bool) {
		for r := 0; r < g.rowLen; r++ {
			for c := 0; c < g.columnLen; c++ {
				if !yield(NewLocation(r, c), g.matrix[r][c]) {
					return
				}
			}
		}
	}
}

// ValidLocation determines if the given location represents a cell that exists
// within the bounds of the grid.
func (g Grid[T]) ValidLocation(l Location) bool {
	return l.Row >= 0 && l.Row < g.rowLen && l.Column >= 0 && l.Column < g.columnLen
}
