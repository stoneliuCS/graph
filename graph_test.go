package graph_test

import (
	"graph"
	"testing"

	"github.com/stretchr/testify/assert"
)

var comparator = func(n1, n2 graph.Node[int]) int {
	return n1.GetVal() - n2.GetVal()
}
var nodeEq = func(n1, n2 graph.Node[int]) bool {
	return n1.GetVal() == n2.GetVal()
}

func TestCreateGraph(t *testing.T) {
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
	assert.NotNil(t, g)
}

func TestAddEdge(t *testing.T) {
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
	n1 := graph.CreateNode(1)
	n2 := graph.CreateNode(2)
	n3 := graph.CreateNode(3)
	g = g.AddEdge(graph.CreateEdge(n1, n2))
	assert.Equal(t, 1, g.GetNumberOfEdges())
	g = g.AddEdge(graph.CreateEdge(n2, n3))
	assert.Equal(t, 2, g.GetNumberOfEdges())
}

func TestDFSGraph1(t *testing.T) {
	// Test case for geeks for geeks
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
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
	dfsTraversal := g.DFS(n0)
	assert.NotEmpty(t, dfsTraversal)
	assert.Equal(t, 5, len(dfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3, n4}, dfsTraversal)
}

func TestDFSGraph2(t *testing.T) {
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
	n0 := graph.CreateNode(0)
	n1 := graph.CreateNode(1)
	n2 := graph.CreateNode(2)
	n3 := graph.CreateNode(3)
	g = g.AddEdge(graph.CreateEdge(n0, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n1))
	g = g.AddEdge(graph.CreateEdge(n0, n3))
	dfsTraversal := g.DFS(n0)
	assert.Equal(t, 4, len(dfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n2, n1, n3}, dfsTraversal)
}

func TestFindNeighboringNodes(t *testing.T) {
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
	n0 := graph.CreateNode(0)
	n1 := graph.CreateNode(1)
	n2 := graph.CreateNode(2)
	n3 := graph.CreateNode(3)
	g = g.AddEdge(graph.CreateEdge(n0, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n1))
	g = g.AddEdge(graph.CreateEdge(n0, n3))
	// Since this is an undirected graph we will find the neighbors for each node
	neighbors := g.FindNeighboringNodes(n0)
	assert.Equal(t, 2, len(neighbors))
	assert.Equal(t, []graph.Node[int]{n2, n3}, neighbors)
	neighbors = g.FindNeighboringNodes(n1)
	assert.Equal(t, 1, len(neighbors))
	assert.Equal(t, []graph.Node[int]{n2}, neighbors)
	neighbors = g.FindNeighboringNodes(n3)
	assert.Equal(t, []graph.Node[int]{n0}, neighbors)
	neighbors = g.FindNeighboringNodes(n2)
	assert.Equal(t, []graph.Node[int]{n0, n1}, neighbors)
}

func TestBFSGraph1(t *testing.T) {
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
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
	bfsTraversal := g.BFS(n0)
	assert.NotEmpty(t, bfsTraversal)
	assert.Equal(t, 5, len(bfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3, n4}, bfsTraversal)
}

func TestBFSGraph2(t *testing.T) {
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
	n0 := graph.CreateNode(0)
	n1 := graph.CreateNode(1)
	n2 := graph.CreateNode(2)
	n3 := graph.CreateNode(3)
	g = g.AddEdge(graph.CreateEdge(n0, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n1))
	g = g.AddEdge(graph.CreateEdge(n0, n3))
	bfsTraversal := g.BFS(n0)
	assert.NotEmpty(t, bfsTraversal)
	assert.Equal(t, 4, len(bfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n2, n3, n1}, bfsTraversal)
}
