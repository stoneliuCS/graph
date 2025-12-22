package graph

import (
	"slices"
)

// Represents a graph of any structure or type.
type Graph[T any] struct {
	edges          []Edge[T]
	nodeComparator func(n1 Node[T], n2 Node[T]) int
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
	val T
}

// Creates a generic empty graph
func Create[T any]() Graph[T] {
	return Graph[T]{
		edges: []Edge[T]{},
	}
}

func (g Graph[T]) AddNodeComparator(comparator func(n1 Node[T], n2 Node[T]) int) Graph[T] {
	return Graph[T]{edges: g.edges, nodeComparator: comparator}
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

// Finds the edges that lead to the given node. Checks using pointer equality.
func (g Graph[T]) FindEdgesThatLeadTo(source Node[T]) []Edge[T] {
	returnEdges := []Edge[T]{}
	for _, e := range g.edges {
		// If e is directed we only check if the edge leads to the source
		if e.directed {
			if &source == &e.v {
				returnEdges = append(returnEdges, e)
			}
		} else {
			if &source == &e.v || &source == &e.u {
				returnEdges = append(returnEdges, e)
			}
		}
	}
	return returnEdges
}

// Finds all edges that lead from the given node. Checks using pointer equality.
func (g Graph[T]) FindEdgesThatLeadFrom(source Node[T]) []Edge[T] {
	returnEdges := []Edge[T]{}
	for _, e := range g.edges {
		// If e is directed we only check if the edge leads to the source
		if e.directed {
			if &source == &e.u {
				returnEdges = append(returnEdges, e)
			}
		} else {
			if &source == &e.v || &source == &e.u {
				returnEdges = append(returnEdges, e)
			}
		}
	}
	return returnEdges
}

// Checks if this graph is directed or undirected. Panics on an empty graph.
func (g Graph[T]) IsDirectedGraph() bool {
	if len(g.edges) == 0 {
		panic("No edges in this graph.")
	}
	allMap := []bool{}
	var firstVal *bool = nil
	for _, e := range g.edges {
		allMap = append(allMap, e.directed)
		if firstVal == nil {
			firstVal = &e.directed
		} else if firstVal != &e.directed {
			panic("Inconsistant Checks")
		}
	}
	return allMap[0]
}

func (g Graph[T]) FindNeighboringNodes(source Node[T]) []Node[T] {
	neighbors := []*Node[T]{}
	// Now check each node which is reachable to this node
	edgesThatLeadFrom := g.FindEdgesThatLeadFrom(source)
	edgesThatLeadTo := g.FindEdgesThatLeadTo(source)
	for _, edge := range edgesThatLeadFrom {
		neighbor := edge.v
		if !slices.Contains(neighbors, &neighbor) {
			neighbors = append(neighbors, &neighbor)
		}
	}
	// Another check, every to edge is also a neighbor
	if g.IsDirectedGraph() {
		for _, edge := range edgesThatLeadTo {
			neighbor := edge.u
			if !slices.Contains(neighbors, &neighbor) {
				neighbors = append(neighbors, &neighbor)
			}
		}
	}
	newNeighbors := []Node[T]{}
	for _, n := range neighbors {
		newNeighbors = append(newNeighbors, *n)
	}
	return newNeighbors
}

// Performs a DFS on this graph from the given source, returns a list of nodes that were visited by DFS in order if this
// graph has a comparator.
func (g Graph[T]) DFS(source Node[T]) []Node[T] {
	var dfsImpl func(src Node[T], visited []*Node[T], acc []Node[T]) []Node[T]
	dfsImpl = func(src Node[T], visited []*Node[T], acc []Node[T]) []Node[T] {
		// Mark the current node as visited
		neighbors := g.FindNeighboringNodes(src)
		if g.nodeComparator != nil {
			slices.SortStableFunc(neighbors, g.nodeComparator)
		}
		for _, neighbor := range neighbors {
			if !slices.Contains(visited, &neighbor) {
				newAcc := append(acc, neighbor)
				newVisited := append(visited, &neighbor)
				return dfsImpl(neighbor, newVisited, newAcc)
			}
		}
		return acc
	}
	return dfsImpl(source, []*Node[T]{&source}, []Node[T]{})
}
