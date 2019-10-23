package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/element"
	"gmidarch/development/messages"
	shared2 "shared"
)

type CalculatorProxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewCalculatorProxy() CalculatorProxy {

	// create a new instance of Server
	r := new(CalculatorProxy)
	r.Behaviour = "B = InvP.e1 -> I_ProcessIn -> InvR.e2 -> TerR.e2 -> I_ProcessOut -> TerP.e1 -> B"

	return *r
}

func (c *CalculatorProxy) Configure(invP, terP, invR, terR *chan messages.SAMessage) CalculatorProxy {

	// configure the state machine
	c.Graph = *graphs.NewExecGraph(6)
	actionChannel := make(chan messages.SAMessage)
	msg := new(messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared2.Invoke, ActionName: "I_Processin", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR, Message: msg}
	c.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR, Message: msg}
	c.Graph.AddEdge(3, 4, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared2.Invoke, ActionName: "I_Processout", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(4, 5, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	c.Graph.AddEdge(5, 0, newEdgeInfo)

	return *c
}

func (CalculatorProxy) I_Processin(msg *messages.SAMessage, info [] *interface{}) {
	inv := shared2.Invocation{}
	inv.Host = "localhost"             // TODO
	inv.Port = shared2.CALCULATOR_PORT // TODO
	inv.Req = msg.Payload.(shared2.Request)

	*msg = messages.SAMessage{Payload: inv}
}

func (CalculatorProxy) I_Processout(msg *messages.SAMessage, info [] *interface{}) {

	result := msg.Payload.([]interface{})
	*msg = messages.SAMessage{Payload: int(result[0].(float64))}
}
