package components

import (
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"shared"
)

type FibonacciInvokerWithMarshaller struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewFibonacciInvokerWithMarshaller() FibonacciInvokerWithMarshaller {

	r := new(FibonacciInvokerWithMarshaller)
	r.Behaviour = "B = InvP.e1 -> I_PrepareToObject -> InvR.e2 -> TerR.e2 -> I_PrepareToSRH -> TerP.e1 -> B"

	return *r
}

func (FibonacciInvokerWithMarshaller) I_PrepareToObject(msg *messages.SAMessage, info [] *interface{}) {

	// unmarshall
	p1 := msg.Payload.([]byte)
	miopPacket := miop.Packet{}
	//err := json.Unmarshal(p1, &miopPacket)    // JSON
	err := msgpack.Unmarshal(p1, &miopPacket)   // Msgpack
	if err != nil {
		log.Fatalf("Marshaller:: Unmarshall:: %s", err)
	}

	//p0 := int(miopPacket.Bd.ReqBody.Body[0].(float64))   // JSON
	p0 := int(miopPacket.Bd.ReqBody.Body[0].(int64))   // Messagepack

	argsTemp := make([]interface{},2)
	argsTemp[0] = p0
	inv := shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	*msg = messages.SAMessage{Payload: inv}
}

func (FibonacciInvokerWithMarshaller) I_PrepareToSRH(msg *messages.SAMessage, info [] *interface{}) {
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
	//r2, err := json.Marshal(miopPacket)     // JSON
	r2, err := msgpack.Marshal(miopPacket)    // Msgpack
	if err != nil {
		log.Fatalf("Marshaller:: Marshall:: %s", err)
	}

	toSRH := make([]interface{}, 1)
	toSRH[0] = r2

	*msg = messages.SAMessage{Payload: toSRH}
}
