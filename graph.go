package graph

import (
	"image/color"
	"os"
	"slices"

	"log"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

// Represents a graph of any structure or type.
type Graph[T any] struct {
	edges          []Edge[T]
	nodeComparator func(n1 Node[T], n2 Node[T]) int
	nodeEqual      func(n1 Node[T], n2 Node[T]) bool
	nodeHash       func(n Node[T]) string
}

// Represents an edge in Graph
type Edge[T any] struct {
	u        Node[T]
	v        Node[T]
	directed bool
	weight   float64
}

// Represents a node in Graph
type Node[T any] struct {
	val T
}

// Creates a generic empty graph with the given comparator and equivalence functions
func CreateGraphFunc[T any](
	comparator func(n1 Node[T], n2 Node[T]) int,
	eqFn func(n1 Node[T], n2 Node[T]) bool,
	hashFn func(n Node[T]) string,
) Graph[T] {
	return Graph[T]{
		edges:          []Edge[T]{},
		nodeComparator: comparator,
		nodeEqual:      eqFn,
		nodeHash:       hashFn,
	}
}

func (g Graph[T]) AddNodeComparator(comparator func(n1 Node[T], n2 Node[T]) int) Graph[T] {
	return Graph[T]{edges: g.edges, nodeComparator: comparator, nodeEqual: g.nodeEqual, nodeHash: g.nodeHash}
}

func (g Graph[T]) AddNodeEqualFn(eqFn func(n1 Node[T], n2 Node[T]) bool) Graph[T] {
	return Graph[T]{edges: g.edges, nodeComparator: g.nodeComparator, nodeEqual: eqFn, nodeHash: g.nodeHash}
}

func (g Graph[T]) AddNodeHash(hash func(n Node[T]) string) Graph[T] {
	return Graph[T]{edges: g.edges, nodeComparator: g.nodeComparator, nodeEqual: g.nodeEqual, nodeHash: hash}
}

// Computes a new graph after adding that edge to this graph. Leaves the original graph unmodified.
func (g Graph[T]) AddEdge(edge Edge[T]) Graph[T] {
	newEdges := append(slices.Clone(g.edges), edge)
	return Graph[T]{
		edges:          newEdges,
		nodeComparator: g.nodeComparator,
		nodeEqual:      g.nodeEqual,
		nodeHash:       g.nodeHash,
	}
}

// Creates an undirected, unweighted edge from u to v.
func CreateEdge[T any](u Node[T], v Node[T]) Edge[T] {
	return Edge[T]{
		u: u,
		v: v,
	}
}

func (e Edge[T]) U() Node[T] {
	return e.u
}

func (e Edge[T]) V() Node[T] {
	return e.v
}

func (e Edge[T]) ToDirected() Edge[T] {
	return Edge[T]{
		u:        e.u,
		v:        e.v,
		directed: true,
		weight:   e.weight,
	}
}

func (e Edge[T]) AddWeight(w float64) Edge[T] {
	return Edge[T]{
		u:        e.u,
		v:        e.v,
		directed: e.directed,
		weight:   w,
	}
}

func CreateNode[T any](val T) Node[T] {
	return Node[T]{
		val: val,
	}
}

func (n Node[T]) GetVal() T {
	return n.val
}

/*
Creates a new graph that interprets all edges as directed. I.e. makes all edges e <-> v to u -> v
*/
func (g Graph[T]) ToDirected() Graph[T] {
	newEdges := []Edge[T]{}
	for _, e := range g.edges {
		newEdges = append(newEdges, e.ToDirected())
	}
	return Graph[T]{
		edges:          newEdges,
		nodeComparator: g.nodeComparator,
		nodeEqual:      g.nodeEqual,
		nodeHash:       g.nodeHash,
	}
}

func (g Graph[T]) GetNumberOfEdges() int {
	return len(g.edges)
}

// Reverses this edge, has no effect on an undirected edge
func (e Edge[T]) Reverse() Edge[T] {
	return Edge[T]{
		u:        e.v,
		v:        e.u,
		directed: e.directed,
		weight:   e.weight,
	}
}

// Finds the edges that lead to the given node. Checks using the given equality function on the graph
func (g Graph[T]) FindEdgesThatLeadTo(source Node[T]) []Edge[T] {
	returnEdges := []Edge[T]{}
	for _, e := range g.edges {
		if e.directed && g.nodeEqual(source, e.v) || (!e.directed && (g.nodeEqual(source, e.v))) {
			returnEdges = append(returnEdges, e)
		} else if !e.directed && g.nodeEqual(source, e.u) {
			// We are going to reverse the direction of the edge
			returnEdges = append(returnEdges, e.Reverse())
		}
	}
	return returnEdges
}

// Finds the edges that lead from the given node. Checks using the given equality function on the graph
func (g Graph[T]) FindEdgesThatLeadFrom(source Node[T]) []Edge[T] {
	returnEdges := []Edge[T]{}
	for _, e := range g.edges {
		if (e.directed && g.nodeEqual(source, e.u)) || (!e.directed && g.nodeEqual(source, e.u)) {
			returnEdges = append(returnEdges, e)
		} else if !e.directed && (g.nodeEqual(source, e.v)) {
			returnEdges = append(returnEdges, e.Reverse())
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
		} else if *firstVal != e.directed {
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

// Performs a DFS on this graph from the given source, returns a list of nodes that were visited by DFS in accordance to
// the graph comparator
func (g Graph[T]) DFS(source Node[T]) []Node[T] {
	var dfsImpl func(src Node[T], visited *[]Node[T], acc *[]Node[T])
	dfsImpl = func(src Node[T], visited *[]Node[T], acc *[]Node[T]) {
		// Mark the current node as visited
		*visited = append(*visited, src)
		neighbors := g.FindNeighboringNodes(src)
		slices.SortStableFunc(neighbors, g.nodeComparator)
		for _, neighbor := range neighbors {
			if !slices.ContainsFunc(*visited, func(n Node[T]) bool { return g.nodeEqual(neighbor, n) }) {
				*acc = append(*acc, neighbor)
				dfsImpl(neighbor, visited, acc)
			}
		}
	}
	acc := []Node[T]{source}
	visited := []Node[T]{}
	dfsImpl(source, &visited, &acc)
	return acc
}

// Performs a BFS on this graph from the given source. Returns a list of nodes that were visited by DFS in accordance to
// the graph comparator
func (g Graph[T]) BFS(source Node[T]) []Node[T] {
	queue := []Node[T]{}
	visited := []Node[T]{}
	bfs := []Node[T]{source}
	queue = append(queue, source)
	// While the queue is not empty
	for len(queue) != 0 {
		// Pop the last neighbor in the queue
		src := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		visited = append(visited, src)
		neighbors := g.FindNeighboringNodes(src)
		slices.SortStableFunc(neighbors, g.nodeComparator)
		for _, neighbor := range neighbors {
			// If we haven't visited this neighbor
			if !slices.ContainsFunc(visited, func(n Node[T]) bool { return g.nodeEqual(neighbor, n) }) {
				bfs = append(bfs, neighbor)
				visited = append(visited, neighbor)
				queue = append(queue, neighbor)
			}
		}
	}
	return bfs
}

// Gets all the nodes in this graph. The nodes are returned as specified in the given node comparator.
func (g Graph[T]) GetNodes() []Node[T] {
	nodes := []Node[T]{}
	for _, edge := range g.edges {
		n1 := edge.u
		n2 := edge.v
		if !slices.ContainsFunc(nodes, func(n Node[T]) bool { return g.nodeEqual(n1, n) }) {
			nodes = append(nodes, n1)
		}
		if !slices.ContainsFunc(nodes, func(n Node[T]) bool { return g.nodeEqual(n2, n) }) {
			nodes = append(nodes, n2)
		}
	}
	slices.SortStableFunc(nodes, g.nodeComparator)
	return nodes
}

// Returns a new graph with all nodes of type U instead of type T. To make the resulting graph valid, one must also pass
// in the corresponding comparator and equivalence functions on that type U
// Maintains the order of the edges from the previous graph.
func MapGraph[T any, U any](
	graph Graph[T],
	mapFn func(Node[T]) Node[U],
	comparator func(n1 Node[U], n2 Node[U]) int,
	eqFn func(n1 Node[U], n2 Node[U]) bool,
	hashFn func(n Node[U]) string,
) Graph[U] {
	newEdges := []Edge[U]{}
	for _, edge := range graph.edges {
		newU := mapFn(edge.u)
		newV := mapFn(edge.v)
		newEdge := Edge[U]{
			u:        newU,
			v:        newV,
			directed: edge.directed,
			weight:   edge.weight,
		}
		newEdges = append(newEdges, newEdge)
	}
	return Graph[U]{
		edges:          newEdges,
		nodeComparator: comparator,
		nodeEqual:      eqFn,
		nodeHash:       hashFn,
	}
}

// Filters the edges from this graph that both nodes on the edge must meet
func FilterGraph[T any](graph Graph[T], filterFn func(Edge[T]) bool) Graph[T] {
	newEdges := []Edge[T]{}
	for _, edge := range graph.edges {
		if filterFn(edge) {
			newEdges = append(newEdges, edge)
		}
	}
	return Graph[T]{
		edges:          newEdges,
		nodeComparator: graph.nodeComparator,
		nodeEqual:      graph.nodeEqual,
		nodeHash:       graph.nodeHash,
	}
}

// Finds the "in-degree" or the number of edges that lead to this source node.
func (g Graph[T]) FindInDegree(source Node[T]) int {
	return len(g.FindEdgesThatLeadTo(source))
}

func (g Graph[T]) FindOutDegree(source Node[T]) int {
	return len(g.FindEdgesThatLeadFrom(source))
}

// Returns all the nodes with indegrees of 0
func (g Graph[T]) GetRootNodes() []Node[T] {
	roots := []Node[T]{}
	for _, node := range g.GetNodes() {
		if g.FindInDegree(node) == 0 {
			roots = append(roots, node)
		}
	}
	return roots
}

// Returns all the nodes with out degrees of 0
func (g Graph[T]) GetLeafNodes() []Node[T] {
	roots := []Node[T]{}
	for _, node := range g.GetNodes() {
		if g.FindOutDegree(node) == 0 {
			roots = append(roots, node)
		}
	}
	return roots
}

// Determines if this graph contains a cycle.
func (g Graph[T]) ContainsCycle() bool {
	var isCyclicImpl func(idx int, node Node[T], visited, stack []bool) bool
	isCyclicImpl = func(idx int, node Node[T], visited, stack []bool) bool {
		if stack[idx] {
			return true
		}
		if visited[idx] {
			return false
		}
		stack[idx] = true
		visited[idx] = true
		for newIdx, neighbor := range g.FindNeighboringNodes(node) {
			if isCyclicImpl(newIdx, neighbor, visited, stack) {
				return true
			}
		}
		stack[idx] = true
		return false
	}
	nodes := g.GetNodes()
	visited := make([]bool, len(nodes))
	stack := make([]bool, len(nodes))
	for idx, node := range nodes {
		if !visited[idx] && isCyclicImpl(idx, node, visited, stack) {
			return true
		}
	}
	return false
}

// Returns a mapping of each nodes neighbors with the supplied hash function used as the key.
func (g Graph[T]) ToAdjacencyMap() map[string][]Node[T] {
	adjMap := map[string][]Node[T]{}
	for _, n := range g.GetNodes() {
		neighbors := g.FindNeighboringNodes(n)
		adjMap[g.nodeHash(n)] = neighbors
	}
	return adjMap
}

// Checks if this graph is a directed acyclic graph
func (g Graph[T]) IsDAG() bool {
	return g.IsDirectedGraph() && !g.ContainsCycle()
}

// Renders an interactive GUI based off the current graph. Useful for visualization and debugging purposes.
func (g Graph[T]) GUI() {
	var run func(*app.Window) error
	run = func(w *app.Window) error {
		theme := material.NewTheme()
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return e.Err
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				title := material.H1(theme, "Hello, Gio")
				maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
				title.Color = maroon
				title.Alignment = text.Middle
				title.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
