package main

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/grid"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	lines, err := input.ReadInput(configuration, "day7")
	if err != nil {
		panic(err)
	}

	diagram := NewDiagram(lines)
	done := false
	for !done {
		shouldContinue := diagram.Cast(true)
		done = !shouldContinue
	}

	fmt.Println(diagram)
	fmt.Println("Total Splits:", diagram.Splits)
	fmt.Println("Total Timelines:", diagram.Timelines())
}

type Cell int

const (
	CellEmpty Cell = iota
	CellStart
	CellSplitter
)

func CellFromRune(r rune) Cell {
	switch r {
	case '.':
		return CellEmpty

	case 'S':
		return CellStart

	case '^':
		return CellSplitter

	default:
		panic(fmt.Errorf("unknown cell value %c", r))
	}
}

func (c Cell) Rune() rune {
	switch c {
	case CellStart:
		return 'S'

	case CellSplitter:
		return '^'

	default:
		return '.'
	}
}

type Beam struct {
	location grid.Location
	count    int
}

func NewBeam(loc grid.Location, count int) *Beam {
	return &Beam{
		location: loc,
		count:    count,
	}
}

type Diagram struct {
	*grid.Grid[Cell]

	Splits int

	activeBeams    []*Beam
	completedBeams []*Beam
}

func NewDiagram(lines []string) *Diagram {
	rowCount := len(lines)
	columnCount := len(lines[0])

	g := grid.NewGrid[Cell](rowCount, columnCount)

	var activeBeams []*Beam
	for r, line := range lines {
		for c, char := range line {
			loc := grid.NewLocation(r, c)
			cell := CellFromRune(char)

			if cell == CellStart {
				activeBeams = append(activeBeams, NewBeam(loc, 1))
			}

			g.SetAt(loc, cell)
		}
	}

	return &Diagram{
		Grid:        g,
		activeBeams: activeBeams,
	}
}

func (d Diagram) beamMap() map[grid.Location]int {
	beamMap := make(map[grid.Location]int)

	for _, beam := range d.activeBeams {
		beamMap[beam.location] = beam.count
	}

	for _, beam := range d.completedBeams {
		beamMap[beam.location] = beam.count
	}

	return beamMap
}

func (d Diagram) Timelines() int {
	var (
		bottomRow = d.RowLen() - 1
		timelines int
	)

	for _, beam := range d.completedBeams {
		if beam.location.Row == bottomRow {
			timelines += beam.count
		}
	}

	return timelines
}

// returns split count and whether or not to continue
func (d *Diagram) Cast(allowOverlap bool) bool {
	consolidatedBeams := make(map[grid.Location]*Beam)
	for _, beam := range d.activeBeams {
		d.completedBeams = append(d.completedBeams, beam)

		newBeams := d.castBeam(beam)
		for _, newBeam := range newBeams {
			if existingBeam, exists := consolidatedBeams[newBeam.location]; exists {
				if allowOverlap {
					existingBeam.count += newBeam.count
				}
			} else {
				consolidatedBeams[newBeam.location] = newBeam
			}
		}
	}

	d.activeBeams = slices.Collect(maps.Values(consolidatedBeams))

	return len(d.activeBeams) != 0
}

func (d *Diagram) castBeam(start *Beam) []*Beam {
	var newBeams []*Beam

	nextLocation := start.location.Down()

	nextCell, ok := d.At(nextLocation)
	if !ok {
		return newBeams
	}

	if nextCell == CellEmpty {
		newBeams = append(newBeams, NewBeam(nextLocation, start.count))
	}

	if nextCell == CellSplitter {
		left := nextLocation.Left()
		if d.ValidLocation(left) {
			newBeams = append(newBeams, NewBeam(left, start.count))
		}

		right := nextLocation.Right()
		if d.ValidLocation(right) {
			newBeams = append(newBeams, NewBeam(right, start.count))
		}

		d.Splits += start.count
	}

	return newBeams
}

func (d Diagram) String() string {
	var (
		row           = 0
		builder       = new(strings.Builder)
		beamLocations = d.beamMap()
	)

	for loc, cell := range d.Iter() {
		if loc.Row != row {
			row = loc.Row
			builder.WriteRune('\n')
		}

		count, _ := beamLocations[loc]
		if cell != CellStart && count > 0 {
			if count != 1 {
				builder.WriteRune('x')
			} else {
				builder.WriteRune('|')
			}
		} else {
			builder.WriteRune(cell.Rune())
		}
	}

	builder.WriteRune('\n')

	builder.WriteRune('\n')
	builder.WriteString("Active Beams: ")
	builder.WriteString(strconv.Itoa(len(d.activeBeams)))
	builder.WriteRune('\n')
	builder.WriteString("Completed Beams: ")
	builder.WriteString(strconv.Itoa(len(d.completedBeams)))
	builder.WriteRune('\n')

	return builder.String()
}
