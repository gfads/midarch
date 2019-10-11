package components

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
)

type Receiver struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewReceiver() Receiver {

	// create a new instance of client
	r := new(Receiver)
	r.Behaviour = "B = InvP.e1 -> I_PrintMessage -> B"

	return *r
}

/*
func (r *Receiver) Configure(invP *chan messages.SAMessage) {

	// configure the state machine
	msg := new(messages.SAMessage)
	r.Graph = *graphs.NewExecGraph(2)
	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)

	actionChannel := make(chan messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Printmessage", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args:args}
	r.Graph.AddEdge(1, 0, newEdgeInfo)

	return
}
*/

func (Receiver) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Receiver:: %v  \n",*msg)
}
