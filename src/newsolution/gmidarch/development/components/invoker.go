package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/element"
	"newsolution/gmidarch/development/messages"
	"newsolution/gmidarch/development/miop"
	"newsolution/shared/shared"
)

type Invoker struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewInvoker() Invoker {

	// create a new instance of Invoker
	r := new(Invoker)
	r.Behaviour = "Invoker = InvP.e1 -> I_DeserialiseMIOP -> InvR.e2 -> TerR.e2 -> I_PrepareToObject -> InvR.e3 -> TerR.e3 -> I_SerialiseMIOP -> InvR.e2 -> TerR.e2 -> I_PrepareToSRH -> TerP.e1 -> Invoker"

	return *r
}

func (i *Invoker) Configure(invP, terP, invR1, terR1, invR2, terR2, invR3, terR3 *chan messages.SAMessage) {

	// configure the state machine
	i.Graph = *graphs.NewExecGraph(12)

	msg := new(messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	i.Graph.AddEdge(0, 1, newEdgeInfo)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_DeserialiseMIOP", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info}
	i.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR1, Message: msg}
	i.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR1, Message: msg}
	i.Graph.AddEdge(3, 4, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Preparetoobject", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info}
	i.Graph.AddEdge(4, 5, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR2, Message: msg}
	i.Graph.AddEdge(5, 6, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR2, Message: msg}
	i.Graph.AddEdge(6, 7, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_SerialiseMIOP", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info}
	i.Graph.AddEdge(7, 8, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR3, Message: msg}
	i.Graph.AddEdge(8, 9, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR3, Message: msg}
	i.Graph.AddEdge(9, 10, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_PreparetoSRH", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info}
	i.Graph.AddEdge(10, 11, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	i.Graph.AddEdge(11, 0, newEdgeInfo)
}

func (Invoker) I_DeserialiseMIOP(msg *messages.SAMessage, info [] *interface{}) {

	argsTemp := make([]interface{}, 1)
	argsTemp[0] = msg.Payload
	msgToMarhsaller := shared.Request{Op: "unmarshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (Invoker) I_PrepareToObject(msg *messages.SAMessage, info [] *interface{}) {
	miopPacket := msg.Payload.(miop.Packet)
	argsTemp := miopPacket.Bd.ReqBody.Body
	inv := shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	*msg = messages.SAMessage{Payload: inv}
}

func (Invoker) I_SerialiseMIOP(msg *messages.SAMessage, info [] *interface{}) {
	r := msg.Payload.(int) // TODO

	// assembly packet
	repHeader := miop.ReplyHeader{Context: "TODO", RequestId: 13, Status: 131313}
	result := make([]interface{}, 1)
	result[0] = r
	repBody := miop.ReplyBody{OperationResult: result}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 2, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{RepHeader: repHeader, RepBody: repBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// configure message
	argsTemp := make([]interface{}, 1)
	argsTemp[0] = miopPacket
	msgToMarhsaller := shared.Request{Op: "marshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (Invoker) I_PrepareToSRH(msg *messages.SAMessage, info [] *interface{}) {
	toSRH := make([]interface{}, 1)
	toSRH[0] = msg.Payload.([]uint8)

	*msg = messages.SAMessage{Payload: toSRH}
}
