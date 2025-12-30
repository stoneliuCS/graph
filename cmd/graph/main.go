package main

import (
	"fmt"
	"graph"
)

func main() {
	g := graph.CreateGraphFunc(
		func(n1, n2 graph.Node[int]) int { return n1.GetVal() - n2.GetVal() },
		func(n1, n2 graph.Node[int]) bool { return n1.GetVal() == n2.GetVal() },
		func(n graph.Node[int]) string { return fmt.Sprint(n.GetVal()) },
	)
	nodes := make([]graph.Node[int], 300)
	for i := range 300 {
		nodes[i] = graph.CreateNode(i + 1)
	}
	for i := 0; i < len(nodes)-1; i++ {
		g = g.AddEdge(graph.CreateEdge(nodes[i], nodes[i+1]))
	}
	g.GUI()
}
