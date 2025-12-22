package main

import (
	"fmt"
	"strings"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/containers"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	lines, err := input.ReadInput(configuration, "day11")
	if err != nil {
		panic(err)
	}

	graph := containers.NewDirectedGraph[string]()
	for _, line := range lines {
		node, outgoing := parseLine(line)
		graph.AddNode(node)
		for _, outgoingNode := range outgoing {
			graph.AddNode(outgoingNode)
			graph.AddEdge(node, outgoingNode, "outgoing")
		}
	}

	// SolvePart1(graph)
	SolvePart2(graph)
}

func SolvePart1(graph *containers.DirectedGraph[string]) {
	paths := FindPathsFromTo("you", "out", graph, nil, nil)

	fmt.Println("Part 1:", len(paths))
}

func SolvePart2(graph *containers.DirectedGraph[string]) {
	count := CountPaths(graph, State{Next: "svr"}, make(map[State]int))

	fmt.Println("Part 2:", count)
}

type State struct {
	Next    string
	SeenDAC bool
	SeenFFT bool
}

func CountPaths(graph *containers.DirectedGraph[string], state State, memo map[State]int) int {
	if count, exists := memo[state]; exists {
		return count
	}

	if state.Next == "out" && state.SeenDAC && state.SeenFFT {
		memo[state] = 1

		return 1
	}

	outgoingNodes, err := graph.GetOutgoingEdges(state.Next, "outgoing")
	if err != nil {
		panic(err)
	}

	var count int
	for outgoingNode := range outgoingNodes {
		nextState := State{
			Next:    outgoingNode,
			SeenDAC: state.SeenDAC || outgoingNode == "dac",
			SeenFFT: state.SeenFFT || outgoingNode == "fft",
		}

		result := CountPaths(graph, nextState, memo)
		count += result
	}

	memo[state] = count

	return count
}

func FindPathsFromTo(from, to string, graph *containers.DirectedGraph[string], pathSeeds []*containers.OrderedSet[string], filter func(string) bool) []*containers.OrderedSet[string] {
	var paths []*containers.OrderedSet[string]
	activePaths := containers.NewHeap(func(a, b *containers.OrderedSet[string]) bool {
		return a.Len() < b.Len()
	})

	start := containers.NewOrderedSet[string]()
	start.Add(from)
	if len(pathSeeds) > 0 {
		for _, pathSeed := range pathSeeds {
			activePaths.Add(pathSeed)
		}
	} else {
		activePaths.Add(start)
	}
	for activePaths.Len() > 0 {
		current, _ := activePaths.Remove()
		finalNode := current.At(-1)
		fmt.Println("Looking at:", current)

		outgoingNodes, err := graph.GetOutgoingEdges(finalNode, "outgoing")
		if err != nil {
			panic(err)
		}

		for outgoingNode := range outgoingNodes {
			if current.Has(outgoingNode) {
				continue
			}

			filtered := false
			if filter != nil {
				filtered = filter(outgoingNode)
			}
			if filtered {
				continue
			}

			newPath := current.Clone()
			newPath.Add(outgoingNode)

			if outgoingNode == to {
				fmt.Println("Found path:", newPath)
				paths = append(paths, newPath)

				continue
			}

			activePaths.Add(newPath)
		}

		fmt.Println("----")
	}

	return paths
}

func parseLine(line string) (string, []string) {
	node, rest, _ := strings.Cut(line, ": ")
	outgoing := strings.Split(rest, " ")

	return node, outgoing

}
