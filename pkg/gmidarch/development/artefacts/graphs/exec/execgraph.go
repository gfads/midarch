package exec

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/shared"
)

type ExecGraph struct {
	NumNodes  int
	ExecEdges [][]ExecEdge
}

type ExecEdge struct {
	From int
	To   int
	Info ExecEdgeInfo
}

type ExecEdgeInfo struct { // TODO define a Action type
	IsInternal     bool // Internal & External
	ActionName     string
	ActionChannel  *chan messages.SAMessage // Channel
	Message        *messages.SAMessage      // Message
	ExternalAction func(*chan messages.SAMessage, *messages.SAMessage)
	InternalAction func(interface{}, []*interface{}, string, *messages.SAMessage, []*interface{}, *bool)
	Info           []*interface{}
}

func NewExecGraph(n int) *ExecGraph {
	return &ExecGraph{
		NumNodes:  n,
		ExecEdges: make([][]ExecEdge, n, shared.NUM_MAX_NODES),
	}
}

func (g *ExecGraph) AddEdge(u, v int, a ExecEdgeInfo) {
	g.ExecEdges[u] = append(g.ExecEdges[u], ExecEdge{From: u, To: v, Info: a})
}

func (g *ExecGraph) AdjacentEdges(u int) []ExecEdge {
	return g.ExecEdges[u]
}
