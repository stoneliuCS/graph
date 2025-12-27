# A Generic Graph Library Written in Go
Algorithms that require graphs are very common. This is a generic graph library that enables many useful features with
its core design philosophy to be functional.

```bash
.
├── go.mod
├── go.sum
├── graph_test.go # graph tests
├── graph.go # Main graph module
└── README.md

1 directory, 5 files
```

## Main Functionalities
The goal of this library is to create a generic graph API with a bunch of useful (and commonly) needed functionalities
these include
- `DFS` _Depth First Search_, enables traversing the graph in a DFS manner.
- `BFS` _Breadth First Search_, enables traversing the graph in a BFS manner.
- `Map` _MapGraph_, enables mapping/translating a graph of type `X` to a graph of type `Y`.
- `Filter` _FilterGraph_, enables filtering of edges on this graph by a specific predicate.
- `Cycle Detection` _ContainsCycle_ determines if a graph contains a cycle.

## API Design
Every single method and function available in `graph` is pure and functional. Meaning that the resulting method
application does not change the underlying graph, instead it returns a new graph underneath. HOWEVER, this does not
gurantee that mutable data types will not be cloned.

## Graphics Support
Under the hood, _graph_ uses _Gio_ to render graph GUIs. If you wish to use the GUI methods available please install
```bash
go install gioui.org/cmd/gogio@latest
```
