package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/element"
	"newsolution/gmidarch/development/messages"
	"newsolution/gmidarch/development/miop"
	"newsolution/shared/shared"
)

type Requestor struct {
	CSP   string
	Graph graphs.ExecGraph
}

func NewRequestor() Requestor {

	// create a new instance of Server
	r := new(Requestor)

	return *r
}

func (r *Requestor) Configure(invP, terP, invR1, terR1, invR2, terR2, invR3, terR3 *chan messages.SAMessage) {

	// configure the state machine
	r.Graph = *graphs.NewExecGraph(12)

	msg := new(messages.SAMessage)
	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)

	info1 := make([]*interface{}, 2) // message, (host+port)
	info1[0] = new(interface{})      // message
	info1[1] = new(interface{})      // host+port
	*info1[0] = msg
	*info1[1] = make([]interface{}, 2)

	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_SerialiseMIOP", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info1}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR1, Message: msg}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR1, Message: msg}
	r.Graph.AddEdge(3, 4, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)

	info2 := make([]*interface{}, 2) // host, port, msg
	info2[0] = new(interface{})
	info2[1] = new(interface{})
	*info2[0] = msg
	*info2[1] = make([] interface{}, 3)
	args2HostPort := make([]interface{}, 2)
	args2HostPort[0] = "localhost"
	args2HostPort[1] = 1313
	*info2[1] = args2HostPort
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_PreparetoCRH", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info2}
	r.Graph.AddEdge(4, 5, newEdgeInfo)

	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR2, Message: msg}
	r.Graph.AddEdge(5, 6, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR2, Message: msg}
	r.Graph.AddEdge(6, 7, newEdgeInfo)

	actionChannel = make(chan messages.SAMessage)
	info3 := make([]*interface{}, 1)
	info3[0] = new(interface{})
	*info3[0] = msg
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_DeserialiseMIOP", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info3}
	r.Graph.AddEdge(7, 8, newEdgeInfo)

	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR3, Message: msg}
	r.Graph.AddEdge(8, 9, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR3, Message: msg}
	r.Graph.AddEdge(9, 10, newEdgeInfo)

	actionChannel = make(chan messages.SAMessage)
	info4 := make([]*interface{}, 1)
	info4[0] = new(interface{})
	*info4[0] = msg
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_PrepareToClient", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info4}
	r.Graph.AddEdge(10, 11, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	r.Graph.AddEdge(11, 0, newEdgeInfo)
}

func (Requestor) I_SerialiseMIOP(msg *messages.SAMessage, info [] *interface{}) { // TODO
	inv := msg.Payload.(shared.Invocation)

	// assembly packet
	reqHeader := miop.RequestHeader{Context: "TODO", RequestId: 13, ResponseExpected: true, Key: 131313, Operation: inv.Req.Op}
	reqBody := miop.RequestBody{Body: inv.Req.Args}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 1, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{ReqHeader: reqHeader, ReqBody: reqBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// store host & port
	hostTemp := new(interface{})
	*hostTemp = inv.Host
	portTemp := new(interface{})
	*portTemp = inv.Port
	*info[0] = hostTemp
	*info[1] = portTemp

	// configure message
	argsTemp := make([]interface{}, 1)
	argsTemp[0] = miopPacket
	msgToMarhsaller := shared.Request{Op: "marshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (Requestor) I_DeserialiseMIOP(msg *messages.SAMessage, info [] *interface{}) {

	argsTemp := make([]interface{}, 1)
	argsTemp[0] = msg.Payload
	msgToMarhsaller := shared.Request{Op: "unmarshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (Requestor) I_PrepareToCRH(msg *messages.SAMessage, info [] *interface{}) {

	hostTemp1:= *info[0]
	hostTemp2 := *hostTemp1.(*interface{})
	hostTemp3 := hostTemp2.(string)

	portTemp1:= *info[1]
	portTemp2 := *portTemp1.(*interface{})
	portTemp3 := portTemp2.(int)

	toCRH := make([]interface{}, 3)
	toCRH[0] = hostTemp3 // host
	toCRH[1] = portTemp3 // port
	toCRH[2] = msg.Payload.([]uint8)

	*msg = messages.SAMessage{Payload: toCRH}
}

func (Requestor) I_PrepareToClient(msg *messages.SAMessage, info [] *interface{}) {
	miopPacket := msg.Payload.(miop.Packet)
	operationResult := miopPacket.Bd.RepBody.OperationResult

	*msg = messages.SAMessage{Payload: operationResult}
}
