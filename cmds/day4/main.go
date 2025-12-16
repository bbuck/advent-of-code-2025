package main

import (
	"fmt"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	scanner, cleanUp, err := input.GetScanner(configuration, "day4")
	if err != nil {
		panic(err)
	}
	defer cleanUp()

	var (
		inputMap []string
		rows     int
		columns  = -1
	)

	for scanner.Scan() {
		rowStr := scanner.Text()

		if columns < 0 {
			columns = len(rowStr)
		}

		rows += 1
		inputMap = append(inputMap, rowStr)
	}

	rollMap := NewMap(rows, columns)
	for row, rowStr := range inputMap {
		for column, item := range rowStr {
			if item != '@' {
				continue
			}

			rollMap.AddRoll(Location{Row: row, Column: column})
		}
	}

	var removedRollCount int
	for {
		var toRemove []Location
		rollMap.ForEach(func(cell Cell, location Location) {
			if cell.Accessible() {
				toRemove = append(toRemove, location)
			}
		})

		for _, l := range toRemove {
			rollMap.RemoveRoll(l)
		}

		removedRollCount += len(toRemove)
		if len(toRemove) == 0 {
			break
		}
	}

	fmt.Println(removedRollCount)
}

type Location struct {
	Row    int
	Column int
}

type Cell struct {
	ContainsRoll bool
	NearbyRolls  int
}

func (c Cell) Accessible() bool {
	return c.ContainsRoll && c.NearbyRolls < 4
}

type Map struct {
	Grid [][]Cell

	rowLen    int
	columnLen int
}

func NewMap(rowCount, columnCount int) *Map {
	var grid [][]Cell
	for range rowCount {
		grid = append(grid, make([]Cell, columnCount))
	}

	return &Map{
		Grid: grid,

		rowLen:    rowCount,
		columnLen: columnCount,
	}
}

func (m Map) ForEach(handler func(Cell, Location)) {
	for r := 0; r < m.RowLen(); r++ {
		for c := 0; c < m.ColumnLen(); c++ {
			handler(m.Grid[r][c], Location{Row: r, Column: c})
		}
	}
}

func (m Map) RowLen() int {
	return m.rowLen
}

func (m Map) ColumnLen() int {
	return m.columnLen
}

func (m *Map) AddRoll(location Location) {
	if !m.ValidLocation(location) {
		return
	}

	cell := m.Grid[location.Row][location.Column]
	cell.ContainsRoll = true
	m.Grid[location.Row][location.Column] = cell

	neighbors := m.Neighbors(location)
	for _, n := range neighbors {
		m.Increment(n)
	}
}

func (m *Map) RemoveRoll(location Location) {
	if !m.ValidLocation(location) {
		return
	}

	cell := m.Grid[location.Row][location.Column]
	cell.ContainsRoll = false
	m.Grid[location.Row][location.Column] = cell

	neighbors := m.Neighbors(location)
	for _, n := range neighbors {
		m.Decrement(n)
	}
}

func (m *Map) Decrement(location Location) {
	m.increment(location, -1)
}

func (m *Map) Increment(location Location) {
	m.increment(location, 1)
}

func (m *Map) increment(location Location, amount int) {
	if !m.ValidLocation(location) {
		return
	}

	cell := m.Grid[location.Row][location.Column]
	cell.NearbyRolls += amount
	m.Grid[location.Row][location.Column] = cell
}

func (m Map) ValidLocation(location Location) bool {
	if location.Row < 0 || location.Row >= m.RowLen() {
		return false
	}

	if location.Column < 0 || location.Column >= m.ColumnLen() {
		return false
	}

	return true
}

func (m Map) Neighbors(location Location) []Location {
	return []Location{
		{Row: location.Row - 1, Column: location.Column - 1},
		{Row: location.Row - 1, Column: location.Column},
		{Row: location.Row - 1, Column: location.Column + 1},
		{Row: location.Row, Column: location.Column - 1},
		{Row: location.Row, Column: location.Column + 1},
		{Row: location.Row + 1, Column: location.Column - 1},
		{Row: location.Row + 1, Column: location.Column},
		{Row: location.Row + 1, Column: location.Column + 1},
	}
}
