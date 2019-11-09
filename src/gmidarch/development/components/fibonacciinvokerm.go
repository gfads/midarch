package components

import (
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"shared"
)

type FibonacciInvokerM struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewFibonacciInvokerM() FibonacciInvokerM {
	r := new(FibonacciInvokerM)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (FibonacciInvokerM) Selector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}) {

	if op == "I_In" {
		return func(msg *messages.SAMessage, info []*interface{}) {
			elem.(FibonacciInvokerM).I_In(msg, info)
		}
	} else { //"I_Out":
		return func(msg *messages.SAMessage, info []*interface{}) {
			elem.(FibonacciInvokerM).I_Out(msg, info)
		}
	}
}

func (FibonacciInvokerM) I_In(msg *messages.SAMessage, info [] *interface{}) {

	// unmarshall
	payload := msg.Payload.([]byte)

	miopPacket := miop.Packet{}
	err := msgpack.Unmarshal(payload, &miopPacket)
	if err != nil {
		log.Fatalf("Fibonacciinvokerwithmarshaller:: %s", err)
	}

	n := miopPacket.Bd.ReqBody.Body[0].(int64)

	// prepare invocation to object
	argsTemp := make([]interface{}, 1, 1)
	argsTemp[0] = int(n)
	inv := shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}

	*msg = messages.SAMessage{Payload: inv}
}

func (FibonacciInvokerM) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	payload := msg.Payload.(int) // TODO - depends on the parameter return

	// assembly packet
	result := make([]interface{}, 1, 1)
	result[0] = payload
	repHeader := miop.ReplyHeader{Context: "TODO", RequestId: 13, Status: 131313}
	repBody := miop.ReplyBody{OperationResult: result}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 2, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{RepHeader: repHeader, RepBody: repBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// configure message
	r, err := msgpack.Marshal(miopPacket)
	if err != nil {
		log.Fatalf("Fibonacciinvokerwithmarshaller:: %s", err)
	}

	toSRH := make([]interface{}, 1, 1)
	toSRH[0] = r

	*msg = messages.SAMessage{Payload: toSRH}
}
