package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"fmt"
)

type Receiver struct {
	CSP      string
	Graph    graphs.GraphExecutable
	InvPChan chan messages.SAMessage
	Msg      messages.SAMessage
}

func NewReceiver(chn *chan messages.SAMessage) Receiver {

	// create a new instance of client
	r := new(Receiver)

	// configure the new instance
	r.CSP = "B = InvP -> I_PrintMessage -> B"
	r.Msg = messages.SAMessage{}
	r.InvPChan = *chn // TODO

	// configure the state machine
	r.Graph = *graphs.NewGraph(2)
	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: Receiver{}.InvP, Message: &r.Msg, ActionChannel: &r.InvPChan, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: Receiver{}.I_PrintMessage, Message: &r.Msg, ActionType: 1, ActionChannel: &actionChannel}
	r.Graph.AddEdge(1, 0, newEdgeInfo)

	return *r
}

func (c Receiver) I_PrintMessage(msg *messages.SAMessage) {
	fmt.Println(*msg)
}

func (c Receiver) InvP(chn *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*chn
}
