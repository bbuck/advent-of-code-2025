package main

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/containers"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	lines, err := input.ReadInput(configuration, "day8")
	if err != nil {
		panic(err)
	}

	var vectors []Vector3
	for _, line := range lines {
		vector, err := ParseVector3(line)
		if err != nil {
			panic(err)
		}

		vectors = append(vectors, vector)
	}

	solvePart2(configuration, vectors)
}

func solvePart2(configuration config.Config, vectors []Vector3) {
	var lines []Line3D
	for i, vector := range vectors {
		for j := i + 1; j < len(vectors); j++ {
			lines = append(lines, NewLine(vector, vectors[j]))
		}
	}

	slices.SortFunc(lines, func(a, b Line3D) int {
		if a.Distance < b.Distance {
			return -1
		}

		if a.Distance > b.Distance {
			return 1
		}

		return 0
	})

	for _, line := range lines {
		fmt.Println(line)
	}
}

// func solvePart1(configuration config.Config, vectors []Vector3) {
// 	graph := containers.NewGraph[Vector3]()
// 	for _, vector := range vectors {
// 		graph.AddNode(vector)
// 	}
//
// 	slices.SortFunc(vectors, func(a, b Vector3) int {
// 		if a.X < b.X {
// 			return -1
// 		}
//
// 		if a.X > b.X {
// 			return 1
// 		}
//
// 		return 0
// 	})
//
// 	targetJunctions := 10
// 	if configuration.Solve {
// 		targetJunctions = 1_000
// 	}
//
// 	heap := containers.NewHeap(func(a, b Line3D) bool {
// 		return a.Distance > b.Distance
// 	})
// 	for i, vector := range vectors {
// 		for j := i + 1; j < len(vectors); j++ {
// 			otherVector := vectors[j]
//
// 			xDist := math.Abs(float64(otherVector.X - vector.X))
// 			maxLine, _ := heap.Peek()
// 			if heap.Len() >= targetJunctions && xDist >= maxLine.Distance {
// 				break
// 			}
//
// 			if heap.Len() < targetJunctions {
// 				line := NewLine(vector, otherVector)
// 				heap.Add(line)
//
// 				continue
// 			}
//
// 			newLine := NewLine(vector, otherVector)
// 			if newLine.Distance < maxLine.Distance {
// 				if line, removed := heap.Remove(); removed {
// 					if (line.Start.X == 162 && line.End.X == 431) || (line.Start.X == 431 && line.End.X == 162) {
// 						fmt.Println("We removed the line in question......")
// 					}
// 				}
//
// 				heap.Add(newLine)
// 			}
// 		}
// 	}
//
// 	line3ds := slices.Collect(heap.Iter())
// 	slices.SortFunc(line3ds, func(a, b Line3D) int {
// 		if a.Distance < b.Distance {
// 			return -1
// 		}
//
// 		if a.Distance > b.Distance {
// 			return 1
// 		}
//
// 		return 0
// 	})
//
// 	for _, line := range line3ds {
// 		graph.AddEdge(line.Start, line.End)
// 	}
//
// 	var (
// 		circuits []Circuit
// 		seen     = containers.NewSet[Vector3]()
// 	)
// 	for vector := range graph.Nodes() {
// 		if seen.Has(vector) {
// 			continue
// 		}
//
// 		circuit := NewCircuit()
// 		buildCircuit(vector, seen, circuit, graph)
//
// 		circuits = append(circuits, circuit)
// 	}
//
// 	slices.SortFunc(circuits, func(a, b Circuit) int {
// 		var (
// 			aLen = a.Len()
// 			bLen = b.Len()
// 		)
//
// 		if aLen < bLen {
// 			return 1
// 		}
//
// 		if aLen > bLen {
// 			return -1
// 		}
//
// 		return 0
// 	})
//
// 	result := circuits[0].Len() * circuits[1].Len() * circuits[2].Len()
//
// 	fmt.Println(result)
// }

func buildCircuit(start Vector3, seen containers.Set[Vector3], circuit Circuit, graph *containers.Graph[Vector3]) {
	if seen.Has(start) {
		return
	}

	seen.Add(start)

	circuit.Add(start)

	edges, _ := graph.GetEdges(start)
	for edgeVector := range edges {
		buildCircuit(edgeVector, seen, circuit, graph)
	}
}

type Vector3 struct {
	X, Y, Z int
}

func ParseVector3(s string) (Vector3, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return Vector3{}, errors.New("invalid string given to ParseVector3")
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return Vector3{}, err
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return Vector3{}, err
	}

	z, err := strconv.Atoi(parts[2])
	if err != nil {
		return Vector3{}, err
	}

	return Vector3{x, y, z}, nil
}

func (v Vector3) DistanceTo(other Vector3) float64 {
	xs := other.X - v.X
	ys := other.Y - v.Y
	zs := other.Z - v.Z

	xs2 := math.Pow(float64(xs), 2)
	ys2 := math.Pow(float64(ys), 2)
	zs2 := math.Pow(float64(zs), 2)

	return math.Sqrt(xs2 + ys2 + zs2)
}

func (v Vector3) String() string {
	return fmt.Sprintf("(%d, %d, %d)", v.X, v.Y, v.Z)
}

type Line3D struct {
	Start, End Vector3
	Distance   float64
}

func NewLine(start, end Vector3) Line3D {
	return Line3D{
		Start:    start,
		End:      end,
		Distance: start.DistanceTo(end),
	}
}

func (l Line3D) String() string {
	return fmt.Sprintf("%s => %s (%f)", l.Start, l.End, l.Distance)
}

type Circuit struct {
	junctions containers.Set[Vector3]
}

func NewCircuit() Circuit {
	return Circuit{
		junctions: containers.NewSet[Vector3](),
	}
}

func (c Circuit) Add(junction Vector3) {
	c.junctions.Add(junction)
}

func (c Circuit) Has(junction Vector3) bool {
	return c.junctions.Has(junction)
}

func (c Circuit) Len() int {
	return len(c.junctions)
}

func (c Circuit) String() string {
	builder := new(strings.Builder)

	builder.WriteString("Circuit (")
	builder.WriteString(strconv.Itoa(len(c.junctions)))
	builder.WriteString(")\n")

	for junction := range c.junctions.Iter() {
		builder.WriteRune('\t')
		builder.WriteString(junction.String())
		builder.WriteRune('\n')
	}

	return builder.String()
}
