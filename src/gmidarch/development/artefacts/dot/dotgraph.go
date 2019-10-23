package dot

type DOTGraph struct {
	NumNodes int
	EdgesDot    [][]DOTEdge
}

type DOTEdge struct {
	From   int
	To     int
	Action string
}

func NewDOTGraph(n int) *DOTGraph {
	return &DOTGraph{
		NumNodes: n,
		EdgesDot:    make([][]DOTEdge, n),
	}
}

func (g *DOTGraph) AddEdge(u, v int, a string) {
	g.EdgesDot[u] = append(g.EdgesDot[u], DOTEdge{From: u, To: v, Action: a})
}

func (g *DOTGraph) AdjacentEdges(u int) []DOTEdge {
	return g.EdgesDot[u]
}
