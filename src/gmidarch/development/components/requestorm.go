package components

import (
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"shared"
)

type RequestorM struct {
	CSP       string
	Graph     graphs.ExecGraph
	Behaviour string
}

func NewRequestorM() RequestorM {

	r := new(RequestorM)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e RequestorM) Selector(elem interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	if op == "I_In" {
		e.I_In(msg, info)
	} else { // "I_Out"
		e.I_Out(msg, info)
	}
}

func (RequestorM) I_In(msg *messages.SAMessage, info [] *interface{}) {
	inv := msg.Payload.(shared.Invocation)

	// assembly packet
	reqHeader := miop.RequestHeader{Context: "TODO", RequestId: 13, ResponseExpected: true, Key: 131313, Operation: inv.Req.Op}
	reqBody := miop.RequestBody{Body: inv.Req.Args}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 1, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{ReqHeader: reqHeader, ReqBody: reqBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// marshall packet
	pckt, err := msgpack.Marshal(miopPacket)
	if err != nil {
		log.Fatalf("Invokerwithmarshaller:: %s", err)
	}

	// store host & port in 'info'
	*info[0] = inv.Host
	*info[1] = inv.Port

	toCRH := make([]interface{}, 3, 3)
	toCRH[0] = inv.Host
	toCRH[1] = inv.Port
	toCRH[2] = pckt

	*msg = messages.SAMessage{Payload: toCRH}
}

func (RequestorM) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	miopPacket := miop.Packet{}

	err := msgpack.Unmarshal(msg.Payload.([]byte), &miopPacket)
	if err != nil {
		log.Fatalf("Requestorwithmarshaller:: %s", err)
	}

	*msg = messages.SAMessage{Payload: miopPacket.Bd.RepBody.OperationResult}
}
