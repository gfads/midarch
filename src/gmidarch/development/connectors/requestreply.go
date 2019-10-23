package connectors

import (
	graphs2 "gmidarch/development/artefacts/graphs"
	element2 "gmidarch/development/element"
	messages2 "gmidarch/development/messages"
)

type Requestreply struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewRequestReply() Requestreply {

	// create a new instance of client
	r := new(Requestreply)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B"

	return *r
}

func (r *Requestreply) Configure (invP, terP, invR, terR *chan messages2.SAMessage) Requestreply {

	// configure the new instance
	//msg := messages.SAMessage{}
	msg := new(messages2.SAMessage)

	// configure the state machine
	r.Graph = *graphs2.NewExecGraph(4)
	newEdgeInfo := graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.InvR, Message: msg, ActionChannel: invR, ActionType: 2}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.TerR, Message: msg, ActionChannel: terR, ActionType: 2}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.TerP, Message: msg, ActionChannel: terP, ActionType: 2}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}