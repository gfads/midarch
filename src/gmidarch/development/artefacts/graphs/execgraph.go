package graphs

import (
	"fmt"
	messages2 "gmidarch/development/messages"
	"os"
	"plugin"
	"reflect"
	"shared/shared"
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

type TypeInternalAction func(any interface{}, name string, msg *messages2.SAMessage, info [] *interface{})
type TypeExternalAction func(*chan messages2.SAMessage, *messages2.SAMessage)

type ExecEdgeInfo struct {
	ActionType     int // Internal & External
	ActionName     string
	ActionChannel  *chan messages2.SAMessage // Channel
	Message        *messages2.SAMessage      // Message
	ExternalAction func(*chan messages2.SAMessage, *messages2.SAMessage)
	InternalAction func(any interface{}, name string, msg *messages2.SAMessage, info [] *interface{})
	Info [] *interface{}
}

func NewExecGraph(n int) *ExecGraph {
	return &ExecGraph{
		NumNodes:  n,
		ExecEdges: make([][]ExecEdge, n),
	}
}

func (g *ExecGraph) AddEdge(u, v int, a ExecEdgeInfo) {
	g.ExecEdges[u] = append(g.ExecEdges[u], ExecEdge{From: u, To: v, Info: a})
}

func (g *ExecGraph) AdjacentEdges(u int) []ExecEdge {
	return g.ExecEdges[u]
}

func (ExecGraph) UpdateGraph(oldGraph ExecGraph, plg plugin.Plugin) ExecGraph {
	newGraph := NewExecGraph(oldGraph.NumNodes)
	*newGraph = oldGraph

	for e1 := range newGraph.ExecEdges {
		for e2 := range newGraph.ExecEdges [e1] {
			action := newGraph.ExecEdges[e1][e2].Info.ActionName
			if shared.IsInternal(action) {
				//fx, err := plg.Lookup(action)
				getType,err := plg.Lookup("Gettype")
				//fx,_ := plg.Lookup("FX")
				if err != nil {
					fmt.Printf("Execgraph:: Old action '%v' is not present in plugin \n", action)
					os.Exit(0)
				}

				//newGraph.ExecEdges[e1][e2].Info.InternalAction = fx
				//oldGraph.ExecEdges[e1][e2].Info.InternalAction = fx

				//tx := fx.(func(any interface{}, name string, msg *messages.SAMessage, info [] *interface{}))
				fmt.Printf("Execgraph:: %v\n", reflect.TypeOf(getType))
				getType.(func() interface{})()
			}
		}
	}

	fx, _ := plg.Lookup("FX")
	fx.(func(int))(3)
	os.Exit(0)
	return *newGraph
}
