package components

import (
	"apps/fibomiddleware/impl"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"os"
)

type FibonacciinvokerM struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewFibonacciInvokerM() FibonacciinvokerM {
	r := new(FibonacciinvokerM)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (e FibonacciinvokerM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Process(msg, info)
}

func (FibonacciinvokerM) I_Process(msg *messages.SAMessage, info [] *interface{}) { // TODO

	// unmarshall
	payload := msg.Payload.([]byte)

	miopPacket := miop.Packet{}
	err := msgpack.Unmarshal(payload, &miopPacket)
	if err != nil {
		log.Fatalf("FibonacciinvokerM:: %s", err)
	}

	//var inv shared.Request
	op := miopPacket.Bd.ReqHeader.Operation
	switch op {
	case "Fibonacci.FiboRPC":
		_p0 := int(miopPacket.Bd.ReqBody.Body[0].(int8)) 	// For better performance on docker
		//var _p0 int 									 	// For general purpose
		//reflectedField := reflect.ValueOf(miopPacket.Bd.ReqBody.Body[0])
		//switch reflectedField.Kind() {
		//	case reflect.Uint16: _p0 = int(miopPacket.Bd.ReqBody.Body[0].(uint16))
		//	case reflect.Uint32: _p0 = int(miopPacket.Bd.ReqBody.Body[0].(uint32))
		//	case reflect.Int8: _p0 = int(miopPacket.Bd.ReqBody.Body[0].(int8))
		//	case reflect.Int64: _p0 = int(miopPacket.Bd.ReqBody.Body[0].(int64))
		//}
		_r := impl.Fibonacci{}.F(_p0)

		// assembly packet
		repHeader := miop.ReplyHeader{Context: "TODO", RequestId: 13, Status: 131313}
		repBody := miop.ReplyBody{OperationResult: _r}
		miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 2, Size: 131313, ByteOrder: true}
		miopBody := miop.Body{RepHeader: repHeader, RepBody: repBody}
		miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

		// configure message
		r, err := msgpack.Marshal(miopPacket)
		if err != nil {
			log.Fatalf("FibonacciinvokerM:: %s", err)
		}

		toSRH := make([]interface{}, 1, 1)
		toSRH[0] = r

		*msg = messages.SAMessage{Payload: toSRH}
	default:
		fmt.Printf("FibonacciInvokerM:: Operation '%v' not implemented by Fibonacci Service\n", op)
		os.Exit(0)
	}
}
