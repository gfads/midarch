package main

import (
"fmt"
"gmidarch/development/framework/messages"
	"gmidarch/development/artefacts/graphs"
	"newsolution/development/element"
	"newsolution/shared"
)

type ClientCalculator struct {
	CSP   string
	Graph graphs.GraphExecutable
	InvP  chan messages.SAMessage
	TerP  chan messages.SAMessage
	InvR  chan messages.SAMessage
	TerR  chan messages.SAMessage
	Msg   messages.SAMessage
}

func GetTypeElem() interface{}{
	return ClientCalculator{}
}

func GetBehaviourExp() string {
	return "B = InvP.e1 -> I_PosInvP -> B"
}

func NewClientCalculator(invR, terR *chan messages.SAMessage) ClientCalculator {

	// create a new instance of client
	r := new(ClientCalculator)

	// configure the new instance
	r.CSP = "B = InvP -> I_ProcessIn -> InvR -> I_ProcessOut -> TerP -> B"
	r.InvR = *invR
	r.TerR = *terR
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(4)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := graphs.EdgeExecutableInfo{InternalAction: ClientCalculator{}.I_SetMessage, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: &r.InvR, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: &r.TerR, Message: &r.Msg}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: ClientCalculator{}.I_PrintMessage, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}

func (ClientCalculator) I_SetMessage(msg *messages.SAMessage) {

	//	if idx < 1 {
	//		idx++
	//t1 = time.Now()
	args := make([]interface{}, 2)
	args[0] = 13
	args[1] = 13
	*msg = messages.SAMessage{Payload: shared.Request{Op: "add", Args: args}}
	//	} else {
	//		os.Exit(0)
	//	}
}

func (ClientCalculator) I_PrintMessage(msg *messages.SAMessage) {

	//fmt.Println(time.Now().Sub(t1))
	fmt.Println(msg.Payload)
}
