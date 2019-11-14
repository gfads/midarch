package graphs

import (
	"gmidarch/development/messages"
	"shared"
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

type TypeInternalAction func(any interface{}, name string, msg *messages.SAMessage, info [] *interface{})
type TypeExternalAction func(*chan messages.SAMessage, *messages.SAMessage)

type ExecEdgeInfo struct { // TODO define a Action type
	IsInternal     bool // Internal & External
	ActionName     string
	ActionChannel  *chan messages.SAMessage // Channel
	Message        *messages.SAMessage      // Message
	ExternalAction func(*chan messages.SAMessage, *messages.SAMessage)
	//InternalAction func(any interface{}, name string, msg *messages.SAMessage, info [] *interface{})
	//InternalAction func(interface{},string) func(*messages.SAMessage,[]*interface{})
	InternalAction func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{})
	Info [] *interface{}
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
