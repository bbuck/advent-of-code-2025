package main

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/containers"
	"bbuck.dev/aoc2025/grid"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	lines, err := input.ReadInput(configuration, "day12")
	if err != nil {
		panic(err)
	}

	presents := make([]Present, 0, 6)
	for range 6 {
		input := lines[0:4]
		present := ParsePresent(input)
		presents = append(presents, present)
		lines = lines[5:]
	}

	var spaces []*Space
	for _, line := range lines {
		space := ParseSpace(line)
		spaces = append(spaces, space)
	}

	SolvePart1(spaces, presents)
}

func SolvePart1(spaces []*Space, presents []Present) {
	var (
		wg     = new(sync.WaitGroup)
		tokens = make(chan struct{}, runtime.NumCPU())
		found  = make(chan struct{}, len(spaces))
	)

	for _, space := range spaces {
		wg.Go(func() {
			fits := space.Fits(presents)

			tokens <- struct{}{}
			defer func() {
				<-tokens
			}()

			if fits {
				found <- struct{}{}
			}

			fmt.Print(".")
		})
	}

	wg.Wait()
	close(found)

	var count int
	for range found {
		count++
	}

	fmt.Println()
	fmt.Println(count)
}

type Space struct {
	Name    string
	Layout  *grid.Grid[bool]
	Counts  []int
	islands *grid.Grid[int]
}

func ParseSpace(line string) *Space {
	var rows, cols int
	counts := make([]int, 6)

	fmt.Sscanf(line, "%dx%d: %d %d %d %d %d %d", &cols, &rows, &counts[0], &counts[1], &counts[2], &counts[3], &counts[4], &counts[5])

	return &Space{
		Name:    fmt.Sprintf("%dx%d", rows, cols),
		Layout:  grid.NewGrid[bool](rows, cols),
		islands: grid.NewGrid[int](rows, cols),
		Counts:  counts,
	}
}

func (space *Space) Islands() []int {
	space.islands.Clear()

	currentIsland := 1
	for loc, set := range space.Layout.Iter() {
		if set {
			continue
		}

		if island, _ := space.islands.At(loc); island == 0 {
			space.checkIsland(loc, currentIsland)

			currentIsland++
		}
	}

	islands := make([]int, currentIsland)
	for _, island := range space.islands.Iter() {
		islands[island-1]++
	}

	return islands
}

func (space *Space) checkIsland(loc grid.Location, currentIsland int) {
	if val, valid := space.islands.At(loc); val != 0 || !valid {
		return
	}

	space.islands.SetAt(loc, currentIsland)

	for _, neighbor := range loc.Neighbors() {
		space.checkIsland(neighbor, currentIsland)
	}
}

func (space *Space) Fits(presents []Present) bool {
	var sum int
	for _, count := range space.Counts {
		sum += count
	}

	expandedPresents := make([]int, 0, sum)
	for i, count := range space.Counts {
		for range count {
			expandedPresents = append(expandedPresents, i)
		}
	}

	slices.SortFunc(expandedPresents, func(a int, b int) int {
		diff := presents[a].Area() - presents[b].Area()

		if diff == 0 {
			return a - b
		}

		return diff
	})

	return space.fitsRecurse(presents, expandedPresents, 0, 0)
}

func (space *Space) fitsRecurse(allPresents []Present, presentIndexes []int, index, startAt int) bool {
	if index >= len(presentIndexes) {
		return true
	}

	presentIndex := presentIndexes[index]
	present := allPresents[presentIndex]

	limit := space.Layout.RowLen() * space.Layout.ColumnLen()
	for locationIndex := startAt; locationIndex < limit; locationIndex++ {
		loc := grid.LocationFromIndex(locationIndex, space.Layout)
		set, _ := space.Layout.At(loc)

		if set {
			continue
		}

		for _, shape := range present.Shapes {
			if shape.PlaceInto(space.Layout, loc) {
				nextStart := 0
				if index+1 < len(presentIndexes) && presentIndexes[index+1] == presentIndex {
					nextStart = locationIndex
				}

				if space.fitsRecurse(allPresents, presentIndexes, index+1, nextStart) {
					return true
				}

				shape.RemoveFrom(space.Layout, loc)
			}
		}
	}

	return false
}

func (space *Space) Debug() string {
	builder := new(strings.Builder)
	row := 0
	for loc, set := range space.Layout.Iter() {
		if loc.Row > row {
			row = loc.Row
			builder.WriteRune('\n')
		}

		if set {
			builder.WriteRune('#')
		} else {
			builder.WriteRune('.')
		}

		builder.WriteRune(' ')
	}
	builder.WriteString("\n------")

	return builder.String()
}

func (space *Space) String() string {
	return fmt.Sprintf("%s: %d %d %d %d %d %d", space.Name, space.Counts[0], space.Counts[1], space.Counts[2], space.Counts[3], space.Counts[4], space.Counts[5])
}

type Shape struct {
	points []grid.Location
	anchor grid.Location
}

func NewShape(points []grid.Location) Shape {
	return Shape{
		points: points,
		anchor: grid.NewLocation(0, 0),
	}
}

func (s Shape) Equals(other Shape) bool {
	set := containers.NewSet[grid.Location]()
	for _, point := range s.points {
		set.Add(point)
	}

	for _, point := range other.points {
		if !set.Has(point) {
			return false
		}
	}

	return true
}

func (s *Shape) Area() int {
	return len(s.points)
}

func (s *Shape) Anchor() {
	g := grid.NewGrid[bool](3, 3)
	s.PlaceInto(g, s.anchor)

	newAnchor := s.anchor
	for loc, set := range g.Iter() {
		if set {
			newAnchor = loc

			break
		}
	}

	if newAnchor == s.anchor {
		return
	}

	for i, point := range s.points {
		s.points[i] = point.Add(s.anchor).Subtract(newAnchor)
	}
	s.anchor = newAnchor
}

func (s Shape) RotateAround(pivot grid.Location) Shape {
	newPoints := make([]grid.Location, len(s.points))
	for i, loc := range s.points {
		rotated := loc.RotateAround(pivot)

		newPoints[i] = rotated
	}

	return NewShape(newPoints)
}

func (s Shape) Mirror() Shape {
	newPoints := make([]grid.Location, len(s.points))
	for i, loc := range s.points {
		var newLoc grid.Location
		switch col := loc.Column; col {
		case 0:
			newLoc = grid.NewLocation(loc.Row, 2)
		case 2:
			newLoc = grid.NewLocation(loc.Row, 0)
		default:
			newLoc = loc
		}

		newPoints[i] = newLoc
	}

	return NewShape(newPoints)
}

func (s Shape) PlaceInto(g *grid.Grid[bool], at grid.Location) bool {
	placed := make([]grid.Location, 0, len(s.points))
	succeeded := true
	for _, point := range s.points {
		translated := point.Add(at)
		if set, valid := g.At(translated); set || !valid {
			succeeded = false

			break
		}

		placed = append(placed, translated)

		g.SetAt(translated, true)
	}

	if !succeeded {
		for _, point := range placed {
			g.SetAt(point, false)
		}
	}

	return succeeded
}

func (s Shape) RemoveFrom(g *grid.Grid[bool], at grid.Location) {
	for _, point := range s.points {
		g.SetAt(point.Add(at), false)
	}
}

func (s Shape) String() string {
	g := grid.NewGrid[bool](3, 3)
	s.PlaceInto(g, s.anchor)

	builder := new(strings.Builder)

	row := 0
	for loc, value := range g.Iter() {
		if loc.Row > row {
			row = loc.Row
			builder.WriteRune('\n')
		}

		if value {
			builder.WriteRune('#')
		} else {
			builder.WriteRune('.')
		}
		builder.WriteRune(' ')
	}

	builder.WriteRune('\n')

	return builder.String()
}

type Present struct {
	Shapes []Shape
}

func ParsePresent(lines []string) Present {
	basePoints := make([]grid.Location, 0, 9)
	for r := 1; r < len(lines); r++ {
		for c, char := range lines[r] {
			if char == '#' {
				basePoints = append(basePoints, grid.NewLocation(r-1, c))
			}
		}
	}

	baseShape := NewShape(basePoints)
	shapes := make([]Shape, 0, 8)
	pivot := grid.NewLocation(1, 1)
	current := baseShape
	for range 4 {
		contains := slices.ContainsFunc(shapes, func(os Shape) bool {
			return current.Equals(os)
		})

		if !contains {
			shapes = append(shapes, current)
		}
		current = current.RotateAround(pivot)
	}

	current = baseShape.Mirror()
	for range 4 {
		contains := slices.ContainsFunc(shapes, func(os Shape) bool {
			return current.Equals(os)
		})

		if !contains {
			shapes = append(shapes, current)
		}
		current = current.RotateAround(pivot)
	}

	for _, shape := range shapes {
		shape.Anchor()
	}

	return Present{shapes}
}

func (p Present) Area() int {
	return p.Shapes[0].Area()
}
