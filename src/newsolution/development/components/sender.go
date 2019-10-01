package components

import (
	"gmidarch/development/framework/messages"
	"newsolution/shared/shared"
	"fmt"
	"newsolution/development/element"
	"newsolution/development/artefacts/exec"
)

type Sender struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewSender() Sender {

	// create a new instance of client
	r := new(Sender)

	return *r

}

func (s *Sender) Configure(invR *chan messages.SAMessage) {

	// Configure state machine
	s.Graph = *exec.NewExecGraph(3)
	actionChannel := make(chan messages.SAMessage)

	msg := new(messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	newEdgeInfo := exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage1", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage2", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage3", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionName: "InvR", ActionType: 2, ActionChannel: invR, Message:msg}
	s.Graph.AddEdge(1, 0, newEdgeInfo)
}

func (Sender) I_Setmessage1(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 1)"}
}
func (Sender) I_Setmessage2(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 2)"}
}
func (Sender) I_Setmessage3(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 3)"}
}

func (Sender) I_Debug(msg *messages.SAMessage) {
	fmt.Printf("Sender:: Debug:: %v \n", msg)
}
