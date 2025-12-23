package graph_test

import (
	"graph"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGraph(t *testing.T) {
	g := graph.Create[int]()
	assert.NotNil(t, g)
}

func TestAddEdge(t *testing.T) {
	g := graph.Create[int]()
	n1 := graph.CreateNode(1)
	n2 := graph.CreateNode(2)
	n3 := graph.CreateNode(3)
	g = g.AddEdge(graph.CreateEdge(n1, n2))
	assert.Equal(t, 1, g.GetNumberOfEdges())
	g = g.AddEdge(graph.CreateEdge(n2, n3))
	assert.Equal(t, 2, g.GetNumberOfEdges())
}

func TestDFS(t *testing.T) {
	// Test case for geeks for geeks
	g := graph.Create[int]()
	n0 := graph.CreateNode(0)
	n1 := graph.CreateNode(1)
	n2 := graph.CreateNode(2)
	n3 := graph.CreateNode(3)
	n4 := graph.CreateNode(4)
	g = g.AddEdge(graph.CreateEdge(n0, n1))
	g = g.AddEdge(graph.CreateEdge(n0, n2))
	g = g.AddEdge(graph.CreateEdge(n1, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n3))
	g = g.AddEdge(graph.CreateEdge(n2, n4))
	g = g.AddNodeComparator(func(n1, n2 graph.Node[int]) int {
		return n1.GetVal() - n2.GetVal()
	})
	g = g.AddNodeEqualFn(func(n1, n2 graph.Node[int]) bool {
		return n1.GetVal() == n2.GetVal()
	})
	dfsTraversal := g.DFS(n0)
	assert.NotEmpty(t, dfsTraversal)
	assert.Equal(t, 5, len(dfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3, n4}, dfsTraversal)
}
