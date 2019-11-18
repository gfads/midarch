package components

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"os"
	"shared"
)

type NaminginvokerM struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewnaminginvokerM() NaminginvokerM {
	r := new(NaminginvokerM)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e NaminginvokerM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { // I_Out
		e.I_Out(msg, info)
	}
}

func (NaminginvokerM) I_In(msg *messages.SAMessage, info [] *interface{}) {

	// unmarshall
	payload := msg.Payload.([]byte)

	miopPacket := miop.Packet{}
	err := msgpack.Unmarshal(payload, &miopPacket)
	if err != nil {
		log.Fatalf("NamingInvokerM:: %s", err)
		os.Exit(0)
	}

	var inv shared.Request
	op := miopPacket.Bd.ReqHeader.Operation
	switch op {
	case "Register":
		_p0 := miopPacket.Bd.ReqBody.Body[0].(string)
		_p1 := miopPacket.Bd.ReqBody.Body[1].(map[string]interface{}) // TODO
		argsTemp := make([]interface{}, 2, 2)
		argsTemp[0] = _p0
		argsTemp[1] = _p1
		inv = shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	case "Lookup":
		fmt.Printf("NamingInvokerM:: Lookup :: HERE\n")
		_p0 := miopPacket.Bd.ReqBody.Body[0].(string)
		argsTemp := make([]interface{}, 1, 1)
		argsTemp[0] = _p0
		inv = shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	case "List":
		argsTemp := make([]interface{}, 0)
		inv = shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	default:
		fmt.Printf("NamingInvokerM:: Operation '%v' not implemented by Naming Service\n",op)
		os.Exit(0)
	}
	*msg = messages.SAMessage{Payload: inv}
}

func (NaminginvokerM) I_Out(msg *messages.SAMessage, info [] *interface{}) { // TODO

	// assembly packet
	repHeader := miop.ReplyHeader{Context: "TODO", RequestId: 13, Status: 131313}
	repBody := miop.ReplyBody{OperationResult: *msg}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 2, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{RepHeader: repHeader, RepBody: repBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// configure message
	r, err := msgpack.Marshal(miopPacket)
	if err != nil {
		log.Fatalf("NamingInvokerM:: %s", err)
	}

	toSRH := make([]interface{}, 1, 1)
	toSRH[0] = r

	*msg = messages.SAMessage{Payload: toSRH}
}