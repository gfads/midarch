package connectors

import (
	graphs2 "gmidarch/development/artefacts/graphs"
	element2 "gmidarch/development/element"
	messages2 "gmidarch/development/messages"
)

type Oneway struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewOneway() Oneway {

	// create a new instance of client
	r := new(Oneway)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> B"

	return *r
}

func (o *Oneway) ConfigureOneWay(invP, invR *chan messages2.SAMessage) {

	// configure the state machine
	//msg := new(messages.SAMessage)
	msg := new(messages2.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	o.Graph = *graphs2.NewExecGraph(2)
	newEdgeInfo := graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	o.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.InvR, Message: msg, ActionChannel: invR, ActionType: 2}
	o.Graph.AddEdge(1, 0, newEdgeInfo)
}