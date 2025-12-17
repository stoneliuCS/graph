package graph

import (
	"slices"

	"github.com/google/uuid"
)

// Represents a graph of any structure or type.
type Graph[T any] struct {
	edges []Edge[T]
}

// Represents an edge in Graph
type Edge[T any] struct {
	u        Node[T]
	v        Node[T]
	directed bool
	weight   *float64
}

// Represents a node in Graph
type Node[T any] struct {
	id  uuid.UUID
	val T
}

// Creates a generic empty graph
func Create[T any]() Graph[T] {
	return Graph[T]{
		edges: []Edge[T]{},
	}
}

// Computes a new graph after adding that edge to this graph. Leaves the original graph unmodified.
func (g Graph[T]) AddEdge(edge Edge[T]) Graph[T] {
	newEdges := append(slices.Clone(g.edges), edge)
	return Graph[T]{
		edges: newEdges,
	}
}

// Creates an undirected, unweighted edge from u to v.
func CreateEdge[T any](u Node[T], v Node[T]) Edge[T] {
	return Edge[T]{
		u: u,
		v: v,
	}
}

func CreateNode[T any](val T) Node[T] {
	return Node[T]{
		val: val,
	}
}

/*
Creates a new graph that interprets all edges as directed.
*/
func (g Graph[T]) ToDirected() Graph[T] {
	newEdges := []Edge[T]{}
	for _, e := range g.edges {
		newEdge := Edge[T]{
			u:        e.u,
			v:        e.v,
			weight:   e.weight,
			directed: true,
		}
		newEdges = append(newEdges, newEdge)
	}
	return Graph[T]{
		edges: newEdges,
	}
}

func (g Graph[T]) GetNumberOfEdges() int {
	return len(g.edges)
}

// Finds the edges that lead to the given node.
func (g Graph[T]) FindEdgesThatLeadTo(source Node[T]) []Edge[T] {
	returnEdges := []Edge[T]{}
	for _, e := range g.edges {
		// If e is directed we only check if the edge leads to the source
		if e.directed {
			if source.id == e.v.id {
				returnEdges = append(returnEdges, e)
			}
		} else {
			if source.id == e.v.id || source.id == e.u.id {
				returnEdges = append(returnEdges, e)
			}
		}
	}
	return returnEdges
}

func (g Graph[T]) FindEdgesThatLeadFrom(source Node[T]) []Edge[T] {
	returnEdges := []Edge[T]{}
	for _, e := range g.edges {
		// If e is directed we only check if the edge leads to the source
		if e.directed {
			if source.id == e.u.id {
				returnEdges = append(returnEdges, e)
			}
		} else {
			if source.id == e.v.id || source.id == e.u.id {
				returnEdges = append(returnEdges, e)
			}
		}
	}
	return returnEdges
}

// Performs a DFS on this graph from the given source, returns a list of nodes that were visited by DFS in order.
// Note that the order of which nodes are visited is arbitrary for DFS.
func (g Graph[T]) DFS(source Node[T]) []Node[T] {
	panic("UnImplemented")
}
