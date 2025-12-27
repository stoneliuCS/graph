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
	g.GUI()
}
