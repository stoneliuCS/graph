package graph_test

import (
	"graph"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper functions
func comparator(n1, n2 graph.Node[int]) int {
	return n1.GetVal() - n2.GetVal()
}
func nodeEq(n1, n2 graph.Node[int]) bool {
	return n1.GetVal() == n2.GetVal()
}
func intToStringNodeMap(n graph.Node[int]) graph.Node[string] {
	val := n.GetVal()
	strVal := strconv.Itoa(val)
	return graph.CreateNode(strVal)
}
func strComparator(n1, n2 graph.Node[string]) int {
	return strings.Compare(n1.GetVal(), n2.GetVal())
}

func strNodeEq(n1, n2 graph.Node[string]) bool {
	return n1.GetVal() == n2.GetVal()
}

func createGraph1() graph.Graph[int] {
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
	return g
}

func createGraph2() graph.Graph[int] {
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
	n0 := graph.CreateNode(0)
	n1 := graph.CreateNode(1)
	n2 := graph.CreateNode(2)
	n3 := graph.CreateNode(3)
	g = g.AddEdge(graph.CreateEdge(n0, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n1))
	g = g.AddEdge(graph.CreateEdge(n0, n3))
	return g
}

func createDirectedCyclicGraph() graph.Graph[int] {
	g := graph.CreateWithEqAndCompFunc(comparator, nodeEq)
	n0 := graph.CreateNode(0)
	n1 := graph.CreateNode(1)
	n2 := graph.CreateNode(2)
	n3 := graph.CreateNode(3)
	g = g.AddEdge(graph.CreateEdge(n0, n1))
	g = g.AddEdge(graph.CreateEdge(n1, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n0))
	g = g.AddEdge(graph.CreateEdge(n2, n3))
	return g.ToDirected()
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
	g := createGraph1()
	n0 := g.GetNodes()[0]
	n1 := g.GetNodes()[1]
	n2 := g.GetNodes()[2]
	n3 := g.GetNodes()[3]
	n4 := g.GetNodes()[4]
	dfsTraversal := g.DFS(n0)
	assert.NotEmpty(t, dfsTraversal)
	assert.Equal(t, 5, len(dfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3, n4}, dfsTraversal)
}

func TestDFSGraph2(t *testing.T) {
	g := createGraph2()
	n0 := g.GetNodes()[0]
	n1 := g.GetNodes()[1]
	n2 := g.GetNodes()[2]
	n3 := g.GetNodes()[3]
	dfsTraversal := g.DFS(n0)
	assert.Equal(t, 4, len(dfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n2, n1, n3}, dfsTraversal)
}

func TestFindNeighboringNodes(t *testing.T) {
	g := createGraph2()
	n0 := g.GetNodes()[0]
	n1 := g.GetNodes()[1]
	n2 := g.GetNodes()[2]
	n3 := g.GetNodes()[3]
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
	g := createGraph1()
	n0 := g.GetNodes()[0]
	n1 := g.GetNodes()[1]
	n2 := g.GetNodes()[2]
	n3 := g.GetNodes()[3]
	n4 := g.GetNodes()[4]
	bfsTraversal := g.BFS(n0)
	assert.NotEmpty(t, bfsTraversal)
	assert.Equal(t, 5, len(bfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3, n4}, bfsTraversal)
}

func TestBFSGraph2(t *testing.T) {
	g := createGraph2()
	n0 := g.GetNodes()[0]
	n1 := g.GetNodes()[1]
	n2 := g.GetNodes()[2]
	n3 := g.GetNodes()[3]
	bfsTraversal := g.BFS(n0)
	assert.NotEmpty(t, bfsTraversal)
	assert.Equal(t, 4, len(bfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n2, n3, n1}, bfsTraversal)
}

func TestGetNodes(t *testing.T) {
	g := createGraph2()
	n0 := g.GetNodes()[0]
	n1 := g.GetNodes()[1]
	n2 := g.GetNodes()[2]
	n3 := g.GetNodes()[3]
	nodes := g.GetNodes()
	assert.Equal(t, 4, len(nodes))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3}, nodes)
}

func TestMapGraph1(t *testing.T) {
	g := createGraph1()
	newG := graph.MapGraph(g, intToStringNodeMap, strComparator, strNodeEq)
	n0 := newG.GetNodes()[0]
	n1 := newG.GetNodes()[1]
	n2 := newG.GetNodes()[2]
	n3 := newG.GetNodes()[3]
	n4 := newG.GetNodes()[4]
	assert.Equal(t, "0", n0.GetVal())
	assert.Equal(t, "1", n1.GetVal())
	assert.Equal(t, "2", n2.GetVal())
	assert.Equal(t, "3", n3.GetVal())
	assert.Equal(t, "4", n4.GetVal())
}

func TestCyclicGraph1(t *testing.T) {
	g := createDirectedCyclicGraph()
	assert.Equal(t, []graph.Node[int]{g.GetNodes()[0], g.GetNodes()[1], g.GetNodes()[2], g.GetNodes()[3]}, g.GetNodes())
	// The indegrees on each node
	assert.Equal(t, 1, g.FindIndegree(g.GetNodes()[0]))
	assert.Equal(t, 1, g.FindIndegree(g.GetNodes()[1]))
	assert.Equal(t, 1, g.FindIndegree(g.GetNodes()[2]))
	assert.Equal(t, 1, g.FindIndegree(g.GetNodes()[3]))
	// The outdegrees on each ndoe
	assert.Equal(t, 1, g.FindOutDegree(g.GetNodes()[0]))
	assert.Equal(t, 1, g.FindOutDegree(g.GetNodes()[1]))
	assert.Equal(t, 2, g.FindOutDegree(g.GetNodes()[2]))
	assert.Equal(t, 0, g.FindOutDegree(g.GetNodes()[3]))
	assert.True(t, g.ContainsCycle())
}
