package main

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	lines, err := input.ReadInput(configuration, "day9")
	if err != nil {
		panic(err)
	}

	var vectors []Vector2
	for _, line := range lines {
		vector, err := ParseVector2(line)
		if err != nil {
			panic(err)
		}

		vectors = append(vectors, vector)
	}

	solvePart2(vectors)
}

// func solvePart1(vectors []Vector2) {
// 	var rects []Rectangle2
// 	for i, vector := range vectors {
// 		for j := i + 1; j < len(vectors); j++ {
// 			rect := NewRectangle2(vector, vectors[j])
//
// 			rects = append(rects, rect)
// 		}
// 	}
//
// 	slices.SortFunc(rects, func(a, b Rectangle2) int {
// 		if a.Area < b.Area {
// 			return 1
// 		}
//
// 		if a.Area > b.Area {
// 			return -1
// 		}
//
// 		return 0
// 	})
//
// 	fmt.Println(rects[0].Area)
// }

func solvePart2(vectors []Vector2) {
	polygon := NewPolygon2(vectors)

	var rects []Rectangle2
	for i, vector := range vectors {
		for j := i + 1; j < len(vectors); j++ {
			rect := NewRectangle2(vector, vectors[j])

			if polygon.ContainsRectangle(rect) {
				rects = append(rects, rect)
			}
		}
	}

	slices.SortFunc(rects, func(a, b Rectangle2) int {
		if a.Area < b.Area {
			return 1
		}

		if a.Area > b.Area {
			return -1
		}

		return 0
	})

	fmt.Println(rects[0])

	fmt.Println(rects[0].Area)
}

type Vector2 struct {
	X, Y int
}

func ParseVector2(s string) (Vector2, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return Vector2{}, errors.New("invalid string given to ParseVector2")
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return Vector2{}, err
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return Vector2{}, err
	}

	return Vector2{x, y}, nil
}

func (v Vector2) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vector2) CrossProduct(other Vector2) int {
	// [ v.X  other.X ]
	// [ v.Y  other.Y ]
	// = (v.X * other.Y) - (other.X * v.Y)

	return (v.X * other.Y) - (other.X * v.Y)
}

type Orientation int8

const (
	OrientationColinear Orientation = iota
	OrientationClockwise
	OrientationCounterClockwise
)

func (o Orientation) String() string {
	switch o {
	case OrientationColinear:
		return "Colinear"

	case OrientationClockwise:
		return "Clockwise"

	default:
		return "Counter-Clockwise"
	}
}

func GetOrientation(a, b, c Vector2) Orientation {
	crossProduct := (b.X-a.X)*(c.Y-a.Y) - (c.X-a.X)*(b.Y-a.Y)

	if crossProduct == 0 {
		return OrientationColinear
	}

	if crossProduct > 0 {
		return OrientationCounterClockwise
	}

	return OrientationClockwise
}

type Line2 struct {
	Start, End Vector2
}

func NewLine2(start, end Vector2) Line2 {
	return Line2{
		Start: start,
		End:   end,
	}
}

func (l Line2) Contains(v Vector2) bool {
	if min(l.Start.X, l.End.X) <= v.X && v.X <= max(l.Start.X, l.End.X) && min(l.Start.Y, l.End.Y) <= v.Y && v.Y <= max(l.Start.Y, l.End.Y) {
		crossProduct := (v.X-l.Start.X)*(l.End.Y-l.Start.Y) - (v.Y-l.Start.Y)*(l.End.X-l.Start.X)

		return crossProduct == 0
	}

	return false
}

func (l Line2) Slope() float32 {
	dy := float32(l.End.Y - l.Start.Y)
	dx := float32(l.End.X - l.Start.X)

	return dy / dx
}

func (l Line2) Intersects(other Line2) bool {
	o1 := GetOrientation(l.Start, l.End, other.Start)
	o2 := GetOrientation(l.Start, l.End, other.End)
	o3 := GetOrientation(other.Start, other.End, l.Start)
	o4 := GetOrientation(other.Start, other.End, l.End)

	if o1 == OrientationColinear || o2 == OrientationColinear || o3 == OrientationColinear || o4 == OrientationColinear {
		return false
	}

	return o1 != o2 && o3 != o4
}

type Polygon2 struct {
	vertices []Vector2
}

func NewPolygon2(vertices []Vector2) Polygon2 {
	return Polygon2{
		vertices: vertices,
	}
}

func (p Polygon2) Segments() []Line2 {
	var segments []Line2

	for i, vertex := range p.vertices {
		next := (i + 1) % len(p.vertices)

		segments = append(segments, NewLine2(vertex, p.vertices[next]))
	}

	return segments
}

func (p Polygon2) ContainsRectangle(rect Rectangle2) bool {
	for vertex := range rect.Vertices() {
		if !p.ContainsVector(vertex) {
			return false
		}
	}

	if !p.ContainsVector(rect.Center()) {
		return false
	}

	polySegments := p.Segments()
	for _, rSegment := range rect.Segments() {
		intersects := slices.ContainsFunc(polySegments, func(pSegment Line2) bool {
			return rSegment.Intersects(pSegment)
		})

		if intersects {
			return false
		}
	}

	return true
}

func (p Polygon2) ContainsVector(vector Vector2) bool {
	intersections := 0

	for i, vertex := range p.vertices {
		next := (i + 1) % len(p.vertices)
		line := NewLine2(vertex, p.vertices[next])

		if line.Contains(vector) {
			return true
		}

		if (line.Start.Y > vector.Y) != (line.End.Y > vector.Y) {
			xIntersect := float64((vector.Y-line.Start.Y)*(line.End.X-line.Start.X))/float64(line.End.Y-line.Start.Y) + float64(line.Start.X)
			if float64(vector.X) <= xIntersect {
				intersections++
			}
		}
	}

	return intersections%2 == 1
}

type Rectangle2 struct {
	LeftCorner, RightCorner Vector2
	Area                    int
}

func NewRectangle2(a, b Vector2) Rectangle2 {
	if b.X < a.X {
		b, a = a, b
	}

	shift := Vector2{min(a.X, b.X), min(a.Y, b.Y)}

	aShifted := a.Sub(shift)
	bShifted := b.Sub(shift)

	length := max(aShifted.X, bShifted.X) + 1
	height := max(aShifted.Y, bShifted.Y) + 1

	return Rectangle2{
		LeftCorner:  a,
		RightCorner: b,
		Area:        length * height,
	}
}

func (r Rectangle2) Center() Vector2 {
	return Vector2{
		X: (r.RightCorner.X + r.LeftCorner.X) / 2,
		Y: (r.RightCorner.Y + r.LeftCorner.Y) / 2,
	}
}

func (r Rectangle2) Vertices() iter.Seq[Vector2] {
	vertices := []Vector2{
		r.LeftCorner,
		{r.LeftCorner.X, r.RightCorner.Y},
		r.RightCorner,
		{r.RightCorner.X, r.LeftCorner.Y},
	}

	return slices.Values(vertices)
}

func (r Rectangle2) Segments() []Line2 {
	vertices := slices.Collect(r.Vertices())

	return []Line2{
		NewLine2(vertices[0], vertices[1]),
		NewLine2(vertices[1], vertices[2]),
		NewLine2(vertices[2], vertices[3]),
		NewLine2(vertices[3], vertices[0]),
	}
}

func (r Rectangle2) String() string {
	return fmt.Sprintf("%s to %s, area = %d", r.LeftCorner, r.RightCorner, r.Area)
}
