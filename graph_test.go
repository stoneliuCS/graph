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
