package containers

import (
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

type EdgeDirection int8

const (
	EdgeDirectionIncoming EdgeDirection = iota
	EdgeDirectionOutgoing
)

type Edges[T comparable] struct {
	Incoming map[string]Set[T]
	Outgoing map[string]Set[T]
}

func newEdges[T comparable]() Edges[T] {
	return Edges[T]{
		Incoming: make(map[string]Set[T]),
		Outgoing: make(map[string]Set[T]),
	}
}

func (e Edges[T]) getRelationshipMap(direction EdgeDirection) map[string]Set[T] {
	if direction == EdgeDirectionOutgoing {
		return e.Outgoing
	}

	return e.Incoming
}

func (e Edges[T]) GetRelationship(direction EdgeDirection, relationship string) Set[T] {
	relationshipMap := e.getRelationshipMap(direction)

	if set, exists := relationshipMap[relationship]; exists {
		return set
	}

	set := NewSet[T]()
	relationshipMap[relationship] = set

	return set
}

func (e Edges[T]) GetRelationships() iter.Seq[string] {
	outgoing := maps.Keys(e.Outgoing)
	incoming := maps.Keys(e.Incoming)

	return func(yield func(string) bool) {
		for relationship := range outgoing {
			if !yield(relationship) {
				return
			}
		}

		for relationship := range incoming {
			if _, exists := e.Outgoing[relationship]; exists {
				continue
			}

			if !yield(relationship) {
				return
			}
		}
	}
}

type DirectedGraph[T comparable] struct {
	graph map[T]Edges[T]
}

func NewDirectedGraph[T comparable]() *DirectedGraph[T] {
	return &DirectedGraph[T]{
		graph: make(map[T]Edges[T]),
	}
}

func (g *DirectedGraph[T]) AddNode(node T) {
	if _, exists := g.graph[node]; !exists {
		g.graph[node] = newEdges[T]()
	}
}

func (g *DirectedGraph[T]) HasNode(node T) bool {
	_, exists := g.graph[node]

	return exists
}

func (g *DirectedGraph[T]) AddEdge(from, to T, relationship string) error {
	if !g.HasNode(from) {
		return errors.New("graph does not have from node")
	}

	if !g.HasNode(to) {
		return errors.New("graph does not have to node")
	}

	g.addEdge(EdgeDirectionOutgoing, from, to, relationship)
	g.addEdge(EdgeDirectionIncoming, to, from, relationship)

	return nil
}

func (g *DirectedGraph[T]) addEdge(direction EdgeDirection, from, to T, relationship string) {
	fromEdges := g.graph[from]

	relationshipSet := fromEdges.GetRelationship(direction, relationship)
	relationshipSet.Add(to)
}

func (g *DirectedGraph[T]) GetEdges(direction EdgeDirection, node T, relationship string) (iter.Seq[T], error) {
	return g.getEdgeIter(direction, node, relationship)
}

func (g *DirectedGraph[T]) GetOutgoingEdges(node T, relationship string) (iter.Seq[T], error) {
	return g.getEdgeIter(EdgeDirectionOutgoing, node, relationship)
}

func (g *DirectedGraph[T]) getEdgeIter(direction EdgeDirection, node T, relationship string) (iter.Seq[T], error) {
	if !g.HasNode(node) {
		return nil, fmt.Errorf("graph does not have node %v", node)
	}

	edges := g.graph[node]
	relationships := edges.getRelationshipMap(direction)

	if nodes, exists := relationships[relationship]; exists {
		return nodes.Iter(), nil
	}

	return slices.Values([]T{}), nil
}

func (g *DirectedGraph[T]) GetIncomingEdges(node T, relationship string) (iter.Seq[T], error) {
	return g.getEdgeIter(EdgeDirectionIncoming, node, relationship)
}

func (g *DirectedGraph[T]) Nodes() iter.Seq[T] {
	return maps.Keys(g.graph)
}

func (g *DirectedGraph[T]) String() string {
	builder := new(strings.Builder)

	for node := range g.Nodes() {
		fmt.Fprint(builder, node)
		builder.WriteRune('\n')

		edges := g.graph[node]
		for relationship := range edges.GetRelationships() {
			builder.WriteString("  ")
			builder.WriteString(relationship)
			builder.WriteRune('\n')

			if outgoingEdges, exists := edges.Outgoing[relationship]; exists {
				for outgoing := range outgoingEdges.Iter() {
					builder.WriteString("   => ")
					fmt.Fprint(builder, outgoing)
					builder.WriteRune('\n')
				}
			}

			if incomingEdges, exists := edges.Incoming[relationship]; exists {
				for incoming := range incomingEdges.Iter() {
					builder.WriteString("   <= ")
					fmt.Fprint(builder, incoming)
					builder.WriteRune('\n')
				}
			}
		}
	}

	return builder.String()
}
