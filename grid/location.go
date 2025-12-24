package grid

// Location represents a position within a 2D grid (multi-dimensional array)
// without formal axes. Locations carry a row and a column data.
type Location struct {
	Row    int
	Column int
}

// NewLocation creates a new location at the given row and column.
func NewLocation(row, column int) Location {
	return Location{
		Row:    row,
		Column: column,
	}
}

func LocationFromIndex[T comparable](index int, g *Grid[T]) Location {
	row := index / g.ColumnLen()
	col := index % g.ColumnLen()

	return NewLocation(row, col)
}

func (l Location) toIndex(columnLen int) int {
	return (l.Row * columnLen) + l.Column
}

// Translate creates a new location with the row and column values shifted.
func (l Location) Translate(rowShift, columnShift int) Location {
	return NewLocation(l.Row+rowShift, l.Column+columnShift)
}

func (l Location) Subtract(other Location) Location {
	return NewLocation(l.Row-other.Row, l.Column-other.Column)
}

func (l Location) Add(other Location) Location {
	return NewLocation(l.Row+other.Row, l.Column+other.Column)
}

func (l Location) RotateAround(pivot Location) Location {
	adjusted := l.Subtract(pivot)
	rotated := NewLocation(adjusted.Column, -adjusted.Row)

	return rotated.Add(pivot)
}

// Up creates a new location that represents the position directly above the
// location in a grid.
func (l Location) Up() Location {
	return l.Translate(-1, 0)
}

// UpLeft creates a new location that represents the position to the top
// left of the current location.
func (l Location) UpLeft() Location {
	return l.Translate(-1, -1)
}

// UpRight creates a new location that represents the position to the top
// right of the current location.
func (l Location) UpRight() Location {
	return l.Translate(-1, 1)
}

// Left creates a new location that represents the position directly to the
// left of the current location.
func (l Location) Left() Location {
	return l.Translate(0, -1)
}

// Right creates a new location that represents the position directly to the
// right of the current location.
func (l Location) Right() Location {
	return l.Translate(0, 1)
}

// Down creates a new location that represents the position directly below the
// current location.
func (l Location) Down() Location {
	return l.Translate(1, 0)
}

// DownLeft creates a new location that represents the position to the bottom
// left of the current location.
func (l Location) DownLeft() Location {
	return l.Translate(1, -1)
}

// DownRight creates a new location that represents the position to the bottom
// right of the current location.
func (l Location) DownRight() Location {
	return l.Translate(1, 1)
}

// Neighbors returns a list of the 8 locations that border the given location
// in the order left, up left, up, up right, right, down right, down, down left
// (clock wise starting from the left).
func (l Location) Neighbors() []Location {
	return []Location{
		l.Left(),
		l.UpLeft(),
		l.Up(),
		l.UpRight(),
		l.Right(),
		l.DownRight(),
		l.Down(),
		l.DownLeft(),
	}
}
