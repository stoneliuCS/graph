package graph_test

import (
	"crypto/sha256"
	"encoding/binary"
	"graph"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TESTING UTILS

type NumberNode struct {
	val int
}

func (n NumberNode) Compare(node graph.Node[int]) int {
	return n.val - node.Val()
}

func (n NumberNode) Equal(node graph.Node[int]) bool {
	return n.val == node.Val()
}

func (n NumberNode) Hash() int {
	return n.val
}

func (n NumberNode) Val() int {
	return n.val
}

type StringNode struct {
	val string
}

func (n StringNode) Compare(node graph.Node[string]) int {
	return strings.Compare(n.val, node.Val())
}

func (n StringNode) Equal(node graph.Node[string]) bool {
	return n.val == node.Val()
}

func (n StringNode) Val() string {
	return n.val
}

func (n StringNode) Hash() int {
	hasher := sha256.New()
	hasher.Write([]byte(n.val))
	hashBytes := hasher.Sum(nil)
	// Take first 8 bytes and convert to int
	return int(binary.BigEndian.Uint64(hashBytes[:8]))
}

// FIXTURES

/*
a ──► b ──► c
▲           │
└───────────┘
*/
func abc() graph.Graph[string] {
	g := graph.CreateDirected[string]()
	g = g.AddEdge(StringNode{"A"}, StringNode{"B"}, 0)
	g = g.AddEdge(StringNode{"B"}, StringNode{"C"}, 0)
	g = g.AddEdge(StringNode{"C"}, StringNode{"A"}, 0)
	return g
}

func abcNoEdges() graph.Graph[string] {
	g := graph.CreateDirected[string]()
	return g.AddNode(StringNode{"A"}).AddNode(StringNode{"B"}).AddNode(StringNode{"C"})
}

// UNIT TESTING

func TestABCGraph(t *testing.T) {
	abc := abc()
	// Basic Tests
	assert.Equal(t, 3, abc.GetNumberOfNodes())
	assert.Equal(t, 3, abc.GetNumberOfEdges())

	// The ABC Graph clearly has a cycle
	assert.True(t, abc.ContainsCycle())

	// DFS Tests
	a_dfs := abc.DFS(StringNode{"A"})
	// Order of DFS should be a, b, c
	assert.Equal(t, []graph.Node[string]{StringNode{"A"}, StringNode{"B"}, StringNode{"C"}}, a_dfs)

	b_dfs := abc.DFS(StringNode{"B"})
	// Order of DFS should be b, c, a
	assert.Equal(t, []graph.Node[string]{StringNode{"B"}, StringNode{"C"}, StringNode{"A"}}, b_dfs)

	c_dfs := abc.DFS(StringNode{"C"})
	// Order of DFS should be b, c, a
	assert.Equal(t, []graph.Node[string]{StringNode{"C"}, StringNode{"A"}, StringNode{"B"}}, c_dfs)

	for _, starting_node := range []graph.Node[string]{StringNode{"A"}, StringNode{"B"}, StringNode{"C"}} {
		dfs := abc.DFS(starting_node)
		bfs := abc.BFS(starting_node)
		assert.Equal(t, dfs, bfs)
	}

	// This graph has a cycle so no possible way to topologically sort this graph.
	assert.Panics(t, func() { abc.GetAllTopologicalSorts() })
}

func TestABCGraphNoEdges(t *testing.T) {
	abc := abcNoEdges()
	topSorts := abc.GetAllTopologicalSorts()
	// There are 3! number of topological sorts since this has no edges.
	assert.Equal(t, 6, len(topSorts))
	// All possible permutations since there are no edges.
	assert.Contains(t, topSorts, []graph.Node[string]{StringNode{"A"}, StringNode{"B"}, StringNode{"C"}})
	assert.Contains(t, topSorts, []graph.Node[string]{StringNode{"A"}, StringNode{"C"}, StringNode{"B"}})
	assert.Contains(t, topSorts, []graph.Node[string]{StringNode{"B"}, StringNode{"C"}, StringNode{"A"}})
	assert.Contains(t, topSorts, []graph.Node[string]{StringNode{"C"}, StringNode{"B"}, StringNode{"A"}})
	assert.Contains(t, topSorts, []graph.Node[string]{StringNode{"B"}, StringNode{"A"}, StringNode{"C"}})
	assert.Contains(t, topSorts, []graph.Node[string]{StringNode{"C"}, StringNode{"A"}, StringNode{"B"}})
}

func TestAdjacencyMaps(t *testing.T) {
	noEdgeAbc := abcNoEdges()
	edgeABC := abc()

	noEdgeAdjEdgeMap := noEdgeAbc.ToAdjacencyEdgeMap()
	noEdgeAdjNodeMap := noEdgeAbc.ToAdjacencyNodeMap()

	adjEdgeMap := edgeABC.ToAdjacencyEdgeMap()
	adjNodeMap := edgeABC.ToAdjacencyNodeMap()
	// There should be an entry for each node in the graph
	assert.Equal(t, 3, len(noEdgeAdjEdgeMap))
	assert.Equal(t, 3, len(noEdgeAdjNodeMap))
	assert.Equal(t, 3, len(adjEdgeMap))
	assert.Equal(t, 3, len(adjNodeMap))

	for _, node := range noEdgeAbc.GetNodes() {
		assert.Empty(t, noEdgeAdjEdgeMap[node.Hash()])
		assert.Empty(t, noEdgeAdjNodeMap[node.Hash()])
	}
}
