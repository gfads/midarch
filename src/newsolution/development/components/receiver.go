package components

import (
	"gmidarch/development/framework/messages"
	"newsolution/development/element"
	"newsolution/shared/shared"
	"fmt"
	"newsolution/development/artefacts/exec"
)

type Receiver struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewReceiver() Receiver {

	// create a new instance of client
	r := new(Receiver)

	return *r
}

func (r *Receiver) Configure(invP *chan messages.SAMessage) {

	// configure the state machine
	msg := new(messages.SAMessage)
	r.Graph = *exec.NewExecGraph(2)
	newEdgeInfo := exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)

	actionChannel := make(chan messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Printmessage", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args:args}
	r.Graph.AddEdge(1, 0, newEdgeInfo)

	return
}

func (Receiver) I_Printmessage(msg *messages.SAMessage) {
	fmt.Printf("Receiver:: %v  \n",*msg)
}
