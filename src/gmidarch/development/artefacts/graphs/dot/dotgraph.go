package dot

import (
	"github.com/gfads/midarch/src/gmidarch/development/connectors"
	"github.com/gfads/midarch/src/gmidarch/development/messages"
)

type DOTGraph struct {
	NumNodes int
	EdgesDot [][]DOTEdge
}

type Action struct {
	Name           string
	IsInternal     bool
	ExternalAction func(*messages.SAMessage, connectors.Connector, string, *interface{}, *bool)
	InternalAction func(interface{}, string, string, *messages.SAMessage, *interface{}, *bool)
	Conn           connectors.Connector // used by external actions only
}

type DOTEdge struct {
	From   int
	To     int
	Action Action
}

func NewDOTGraph(n int) *DOTGraph {
	return &DOTGraph{
		NumNodes: n,
		EdgesDot: make([][]DOTEdge, n),
	}
}

func (g *DOTGraph) AddEdge(u, v int, a string) {
	action := Action{}
	action.Name = a
	g.EdgesDot[u] = append(g.EdgesDot[u], DOTEdge{From: u, To: v, Action: action})
}

func (g *DOTGraph) AdjacentEdges(u int) []DOTEdge {
	return g.EdgesDot[u]
}
