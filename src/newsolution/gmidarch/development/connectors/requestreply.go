package connectors

import (
	"newsolution/gmidarch/development/messages"
	"newsolution/gmidarch/development/element"
	"newsolution/gmidarch/development/artefacts/graphs"
)

type Requestreply struct {
	Behaviour   string
	Graph graphs.ExecGraph
}

func NewRequestReply() Requestreply {

	// create a new instance of client
	r := new(Requestreply)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B"

	return *r
}

func (r *Requestreply) Configure (invP, terP, invR, terR *chan messages.SAMessage) Requestreply {

	// configure the new instance
	//msg := messages.SAMessage{}
	msg := new(messages.SAMessage)

	// configure the state machine
	r.Graph = *graphs.NewExecGraph(4)
	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, Message: msg, ActionChannel: invR, ActionType: 2}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, Message: msg, ActionChannel: terR, ActionType: 2}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, Message: msg, ActionChannel: terP, ActionType: 2}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}