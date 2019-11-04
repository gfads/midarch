package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"shared"
)

type Requestor struct {
	CSP       string
	Graph     graphs.ExecGraph
	Behaviour string
}

func NewRequestor() Requestor {

	r := new(Requestor)
	r.Behaviour = "B = InvP.e1 -> I_SerialiseMIOP -> InvR.e2 -> TerR.e2 -> I_PrepareToCRH -> InvR.e3 -> TerR.e3 -> I_DeserialiseMIOP -> InvR.e2 -> TerR.e2 -> I_PrepareToClient -> TerP.e1 -> B"

	return *r
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

func (Requestor) I_DeserialiseMIOP(msg *messages.SAMessage, info [] *interface{}) {

	argsTemp := make([]interface{}, 1)
	argsTemp[0] = msg.Payload
	msgToMarhsaller := shared.Request{Op: "unmarshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (Requestor) I_PrepareToClient(msg *messages.SAMessage, info [] *interface{}) {
	miopPacket := msg.Payload.(miop.Packet)
	operationResult := miopPacket.Bd.RepBody.OperationResult

	*msg = messages.SAMessage{Payload: operationResult}
}