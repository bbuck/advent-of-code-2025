package containers

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"strings"
)

type graphable interface {
	comparable
	fmt.Stringer
}

type Graph[T graphable] struct {
	graph map[T]Set[T]
}

func NewGraph[T graphable]() *Graph[T] {
	return &Graph[T]{
		graph: make(map[T]Set[T]),
	}
}

func (g *Graph[T]) AddNode(node T) {
	if _, exists := g.graph[node]; !exists {
		g.graph[node] = NewSet[T]()
	}
}

func (g *Graph[T]) HasNode(node T) bool {
	_, exists := g.graph[node]

	return exists
}

func (g *Graph[T]) AddEdge(from, to T) error {
	if !g.HasNode(from) {
		return errors.New("graph does not have from node")
	}

	if !g.HasNode(to) {
		return errors.New("graph does not have to node")
	}

	fromEdges := g.graph[from]
	if !fromEdges.Has(to) {
		fromEdges.Add(to)
		g.AddEdge(to, from)
	}

	return nil
}

func (g *Graph[T]) GetEdges(node T) (iter.Seq[T], error) {
	if !g.HasNode(node) {
		return nil, errors.New("graph does not have node")
	}

	edges := g.graph[node]

	return edges.Iter(), nil
}

func (g *Graph[T]) Nodes() iter.Seq[T] {
	return maps.Keys(g.graph)
}

func (g *Graph[T]) Inspect() string {
	builder := new(strings.Builder)

	for node := range g.Nodes() {
		builder.WriteString(node.String())
		builder.WriteRune('\n')

		edges, _ := g.GetEdges(node)
		for edge := range edges {
			builder.WriteString(" => ")
			builder.WriteString(edge.String())
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}
