package components

import (
	"encoding/json"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"shared"
)

type RequestorWithMarhsaller struct {
	CSP       string
	Graph     graphs.ExecGraph
	Behaviour string
}

func NewRequestorWithMarhsaller() RequestorWithMarhsaller {

	r := new(RequestorWithMarhsaller)
	r.Behaviour = "B = InvP.e1 -> I_PrepareToCRH -> InvR.e2 -> TerR.e2 -> I_PrepareToClient -> TerP.e1 -> B"

	return *r
}

func (RequestorWithMarhsaller) I_SerialiseMIOP(msg *messages.SAMessage, info [] *interface{}) { // TODO
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

func (RequestorWithMarhsaller) I_PrepareToCRH(msg *messages.SAMessage, info [] *interface{}) {
	inv := msg.Payload.(shared.Invocation)

	// assembly packet
	reqHeader := miop.RequestHeader{Context: "TODO", RequestId: 13, ResponseExpected: true, Key: 131313, Operation: inv.Req.Op}
	reqBody := miop.RequestBody{Body: inv.Req.Args}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 1, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{ReqHeader: reqHeader, ReqBody: reqBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// marshall packet
	miopPacketSerialised, err := json.Marshal(miopPacket)
	if err != nil {
		log.Fatalf("Marshaller:: Marshall:: %s", err)
	}

	// store host & port
	hostTemp := new(interface{})
	*hostTemp = inv.Host
	portTemp := new(interface{})
	*portTemp = inv.Port
	*info[0] = hostTemp
	*info[1] = portTemp

	toCRH := make([]interface{}, 3)
	toCRH[0] = inv.Host
	toCRH[1] = inv.Port
	toCRH[2] = miopPacketSerialised

	*msg = messages.SAMessage{Payload: toCRH}
}

func (RequestorWithMarhsaller) I_DeserialiseMIOP(msg *messages.SAMessage, info [] *interface{}) {

	argsTemp := make([]interface{}, 1)
	argsTemp[0] = msg.Payload
	msgToMarhsaller := shared.Request{Op: "unmarshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (RequestorWithMarhsaller) I_PrepareToClient(msg *messages.SAMessage, info [] *interface{}) {
	miopPacket := miop.Packet{}
	err := json.Unmarshal(msg.Payload.([]byte), &miopPacket)
	if err != nil {
		log.Fatalf("Marshaller:: Unmarshall:: %s", err)
	}
	*msg = messages.SAMessage{Payload: miopPacket.Bd.RepBody.OperationResult}
}