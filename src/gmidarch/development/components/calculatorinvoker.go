package components

import (
	graphs2 "gmidarch/development/artefacts/graphs"
	messages2 "gmidarch/development/messages"
	miop2 "gmidarch/development/miop"
	"shared/shared"
)

type Calculatorinvoker struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewCalculatorinvoker() Invoker {

	// create a new instance of Invoker
	r := new(Invoker)
	r.Behaviour = "B = InvP.e1 -> I_DeserialiseMIOP -> InvR.e2 -> TerR.e2 -> I_PrepareToObject -> InvR.e3 -> TerR.e3 -> I_SerialiseMIOP -> InvR.e2 -> TerR.e2 -> I_PrepareToSRH -> TerP.e1 -> B"

	return *r
}

/*
func (i *Calculatorinvoker) Configure(invP, terP, invR1, terR1, invR2, terR2, invR3, terR3 *chan messages.SAMessage) {

	// configure the state machine
	i.Graph = *graphs.NewExecGraph(12)

	msg := new(messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	i.Graph.AddEdge(0, 1, newEdgeInfo)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_DeserialiseMIOP", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args:args}
	i.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR1, Message: msg}
	i.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR1, Message: msg}
	i.Graph.AddEdge(3, 4, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_PrepareToObject", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args:args}
	i.Graph.AddEdge(4, 5, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR2, Message: msg}
	i.Graph.AddEdge(5, 6, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR2, Message: msg}
	i.Graph.AddEdge(6, 7, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_SerialiseMIOP", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args: args}
	i.Graph.AddEdge(7, 8, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR3, Message: msg}
	i.Graph.AddEdge(8, 9, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR3, Message: msg}
	i.Graph.AddEdge(9, 10, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_PrepareToSRH", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args: args}
	i.Graph.AddEdge(10, 11, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	i.Graph.AddEdge(11, 0, newEdgeInfo)
}
*/
func (Calculatorinvoker) I_DeserialiseMIOP(msg *messages2.SAMessage, info [] *interface{}) {

	argsTemp := make([]interface{}, 1)
	argsTemp[0] = msg.Payload
	msgToMarhsaller := shared.Request{Op: "unmarshall", Args: argsTemp}

	*msg = messages2.SAMessage{Payload: msgToMarhsaller}
}

func (Calculatorinvoker) I_PrepareToObject(msg *messages2.SAMessage) {
	miopPacket := msg.Payload.(miop2.Packet)
	argsTemp := miopPacket.Bd.ReqBody.Body
	inv := shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	*msg = messages2.SAMessage{Payload: inv}
}

func (Calculatorinvoker) I_SerialiseMIOP(msg *messages2.SAMessage, info [] *interface{}) {
	r := msg.Payload.(int) // TODO

	// assembly packet
	repHeader := miop2.ReplyHeader{Context: "TODO", RequestId: 13, Status: 131313}
	result := make([]interface{}, 1)
	result[0] = r
	repBody := miop2.ReplyBody{OperationResult: result}
	miopHeader := miop2.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 2, Size: 131313, ByteOrder: true}
	miopBody := miop2.Body{RepHeader: repHeader, RepBody: repBody}
	miopPacket := miop2.Packet{Hdr: miopHeader, Bd: miopBody}

	// configure message
	argsTemp := make([]interface{}, 1)
	argsTemp[0] = miopPacket
	msgToMarhsaller := shared.Request{Op: "marshall", Args: argsTemp}

	*msg = messages2.SAMessage{Payload: msgToMarhsaller}
}

func (Calculatorinvoker) I_PrepareToSRH(msg *messages2.SAMessage, info [] *interface{}) {
	toSRH := make([]interface{}, 1)
	toSRH[0] = msg.Payload.([]uint8)

	*msg = messages2.SAMessage{Payload: toSRH}
}
