package components

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
)

type Sender struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewSender() Sender {

	// create a new instance of client
	r := new(Sender)
	r.Behaviour = "B = I_Setmessage1 -> InvR.e1 -> B [] I_Setmessage2 -> InvR.e1 -> B [] I_Setmessage3 -> InvR.e1 -> B"

	return *r

}

/*
func (s *Sender) Configure(invR *chan messages.SAMessage) {

	// Configure state machine
	s.Graph = *graphs.NewExecGraph(3)
	actionChannel := make(chan messages.SAMessage)

	msg := new(messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	newEdgeInfo := graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage1", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage2", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage3", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionName: "InvR", ActionType: 2, ActionChannel: invR, Message:msg}
	s.Graph.AddEdge(1, 0, newEdgeInfo)
}
*/

func (Sender) I_Setmessage1(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 1)"}
}
func (Sender) I_Setmessage2(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 2)"}
}
func (Sender) I_Setmessage3(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 3)"}
}
func (Sender) I_Debug(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Sender:: Debug:: %v \n", msg)
}
