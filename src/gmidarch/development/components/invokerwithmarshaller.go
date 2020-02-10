package components

import (
	"encoding/json"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"shared"
)

type InvokerWithMarshaller struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewInvokerWithMarshaller() InvokerWithMarshaller {

	r := new(InvokerWithMarshaller)
	r.Behaviour = "B = InvP.e1 -> I_PrepareToObject -> InvR.e2 -> TerR.e2 -> I_PrepareToSRH -> TerP.e1 -> B"

	return *r
}

func (InvokerWithMarshaller) I_PrepareToObject(msg *messages.SAMessage, info [] *interface{}) {
	payload := msg.Payload.([]byte)

	pckt := miop.Packet{}
	err := json.Unmarshal(payload, &pckt)
	if err != nil {
		log.Fatalf("Invokerwithmarshaller:: %s", err)
	}

 	//miopPacket := pckt
	p0 := int(pckt.Bd.ReqBody.Body[0].(float64))   // JSON
	//p0 := int(pckt.Bd.ReqBody.Body[0].(int64))       // Messagepack

	argsTemp := make([]interface{},1)
	argsTemp[0] = p0
	inv := shared.Request{Op: pckt.Bd.ReqHeader.Operation, Args: argsTemp}
	*msg = messages.SAMessage{Payload: inv}
}

func (InvokerWithMarshaller) I_PrepareToSRH(msg *messages.SAMessage, info [] *interface{}) {
	r := msg.Payload.(int) // TODO

	// assembly packet
	repHeader := miop.ReplyHeader{Context: "TODO", RequestId: 13, Status: 131313}
	result := make([]interface{}, 1)
	result[0] = r
	repBody := miop.ReplyBody{OperationResult: result}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 2, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{RepHeader: repHeader, RepBody: repBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

    p1 := miopPacket
	r2, err := json.Marshal(p1)
	if err != nil {
		log.Fatalf("Invokerwithmarshaller:: %s", err)
	}

	toSRH := make([]interface{}, 1)
	toSRH[0] = r2

	*msg = messages.SAMessage{Payload: toSRH}
}
