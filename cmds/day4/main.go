package main

import (
	"fmt"
	"iter"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/grid"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	lines, err := input.ReadInput(configuration, "day4")
	if err != nil {
		panic(err)
	}

	rollMap := NewMap(len(lines), len(lines[0]))
	for row, rowStr := range lines {
		for column, item := range rowStr {
			if item != '@' {
				continue
			}

			rollMap.AddRoll(grid.NewLocation(row, column))
		}
	}

	var removedRollCount int
	for {
		var toRemove []grid.Location
		for location, cell := range rollMap.Iter() {
			if cell.Accessible() {
				toRemove = append(toRemove, location)
			}
		}

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

type Cell struct {
	ContainsRoll bool
	NearbyRolls  int
}

func (c Cell) Accessible() bool {
	return c.ContainsRoll && c.NearbyRolls < 4
}

type Map struct {
	grid *grid.Grid[Cell]
}

func NewMap(rowCount, columnCount int) *Map {
	return &Map{
		grid: grid.NewGrid[Cell](rowCount, columnCount),
	}
}

func (m Map) Iter() iter.Seq2[grid.Location, Cell] {
	return m.grid.Iter()
}

func (m Map) AddRoll(location grid.Location) {
	if !m.grid.ValidLocation(location) {
		return
	}

	m.grid.UpdateAt(location, func(cell Cell) Cell {
		cell.ContainsRoll = true

		return cell
	})

	for _, n := range location.Neighbors() {
		m.Increment(n)
	}
}

func (m Map) RemoveRoll(location grid.Location) {
	if !m.grid.ValidLocation(location) {
		return
	}

	m.grid.UpdateAt(location, func(cell Cell) Cell {
		cell.ContainsRoll = false

		return cell
	})

	for _, n := range location.Neighbors() {
		m.Decrement(n)
	}
}

func (m Map) Decrement(location grid.Location) {
	m.grid.UpdateAt(location, func(cell Cell) Cell {
		cell.NearbyRolls--

		return cell
	})
}

func (m Map) Increment(location grid.Location) {
	m.grid.UpdateAt(location, func(cell Cell) Cell {
		cell.NearbyRolls++

		return cell
	})
}
