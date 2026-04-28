package graph_test

import (
	"crypto/sha256"
	"encoding/binary"
	"graph"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

var n0 graph.Node[int] = NumberNode{0}
var n1 graph.Node[int] = NumberNode{1}
var n2 graph.Node[int] = NumberNode{2}
var n3 graph.Node[int] = NumberNode{3}
var n4 graph.Node[int] = NumberNode{4}

func intToStringNodeMap(n graph.Node[int]) graph.Node[string] {
	val := n.Val()
	strVal := strconv.Itoa(val) // int → string
	return StringNode{strVal}
}

func createGraph1() graph.Graph[int] {
	g := graph.CreateGraphFunc[int]()
	g = g.AddEdge(graph.CreateEdge(n0, n1))
	g = g.AddEdge(graph.CreateEdge(n0, n2))
	g = g.AddEdge(graph.CreateEdge(n1, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n3))
	g = g.AddEdge(graph.CreateEdge(n2, n4))
	return g
}

func createGraph2() graph.Graph[int] {
	g := graph.CreateGraphFunc[int]()
	g = g.AddEdge(graph.CreateEdge(n0, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n1))
	g = g.AddEdge(graph.CreateEdge(n0, n3))
	return g
}

func createDirectedCyclicGraph() graph.Graph[int] {
	g := graph.CreateGraphFunc[int]()
	n0 := NumberNode{0}
	n1 := NumberNode{1}
	n2 := NumberNode{2}
	n3 := NumberNode{3}
	g = g.AddEdge(graph.CreateEdge(n0, n1))
	g = g.AddEdge(graph.CreateEdge(n1, n2))
	g = g.AddEdge(graph.CreateEdge(n2, n0))
	g = g.AddEdge(graph.CreateEdge(n2, n3))
	return g.ToDirected()
}

func TestCreateGraph(t *testing.T) {
	g := graph.CreateGraphFunc[int]()
	assert.NotNil(t, g)
}

func TestAddEdge(t *testing.T) {
	g := graph.CreateGraphFunc[int]()
	g = g.AddEdge(graph.CreateEdge(n1, n2))
	assert.Equal(t, 1, g.GetNumberOfEdges())
	g = g.AddEdge(graph.CreateEdge(n2, n3))
	assert.Equal(t, 2, g.GetNumberOfEdges())
}

func TestDFSGraph1(t *testing.T) {
	// Test case for geeks for geeks
	g := createGraph1()
	dfsTraversal := g.DFS(n0)
	assert.NotEmpty(t, dfsTraversal)
	assert.Equal(t, 5, len(dfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3, n4}, dfsTraversal)
}

func TestDFSGraph2(t *testing.T) {
	g := createGraph2()
	dfsTraversal := g.DFS(n0)
	assert.Equal(t, 4, len(dfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n2, n1, n3}, dfsTraversal)
}

func TestFindNeighboringNodes(t *testing.T) {
	g := createGraph2()
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
	bfsTraversal := g.BFS(n0)
	assert.NotEmpty(t, bfsTraversal)
	assert.Equal(t, 5, len(bfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3, n4}, bfsTraversal)
}

func TestBFSGraph2(t *testing.T) {
	g := createGraph2()
	bfsTraversal := g.BFS(n0)
	assert.NotEmpty(t, bfsTraversal)
	assert.Equal(t, 4, len(bfsTraversal))
	assert.Equal(t, []graph.Node[int]{n0, n2, n3, n1}, bfsTraversal)
}

func TestGetNodes(t *testing.T) {
	g := createGraph2()
	nodes := g.GetNodes()
	assert.Equal(t, 4, len(nodes))
	assert.Equal(t, []graph.Node[int]{n0, n1, n2, n3}, nodes)
}

func TestMapGraph1(t *testing.T) {
	g := createGraph1()
	newG := graph.MapGraph(g, intToStringNodeMap)
	n0 := newG.GetNodes()[0]
	n1 := newG.GetNodes()[1]
	n2 := newG.GetNodes()[2]
	n3 := newG.GetNodes()[3]
	n4 := newG.GetNodes()[4]
	assert.Equal(t, "0", n0.Val())
	assert.Equal(t, "1", n1.Val())
	assert.Equal(t, "2", n2.Val())
	assert.Equal(t, "3", n3.Val())
	assert.Equal(t, "4", n4.Val())
}

func TestCyclicGraph1(t *testing.T) {
	g := createDirectedCyclicGraph()
	assert.Equal(t, []graph.Node[int]{g.GetNodes()[0], g.GetNodes()[1], g.GetNodes()[2], g.GetNodes()[3]}, g.GetNodes())
	// The indegrees on each node
	assert.Equal(t, 1, g.FindInDegree(g.GetNodes()[0]))
	assert.Equal(t, 1, g.FindInDegree(g.GetNodes()[1]))
	assert.Equal(t, 1, g.FindInDegree(g.GetNodes()[2]))
	assert.Equal(t, 1, g.FindInDegree(g.GetNodes()[3]))
	// The outdegrees on each node
	assert.Equal(t, 1, g.FindOutDegree(g.GetNodes()[0]))
	assert.Equal(t, 1, g.FindOutDegree(g.GetNodes()[1]))
	assert.Equal(t, 2, g.FindOutDegree(g.GetNodes()[2]))
	assert.Equal(t, 0, g.FindOutDegree(g.GetNodes()[3]))
	assert.True(t, g.ContainsCycle())
}

func TestReverseEdge(t *testing.T) {
	edge := graph.CreateEdge(n0, n1)
	reversedEdge := edge.Reverse()
	assert.Equal(t, edge.U(), reversedEdge.V())
	assert.Equal(t, edge.V(), reversedEdge.U())
}
