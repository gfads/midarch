package main

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

func Gettype() interface{} {
	return FibonacciinvokerM{}
}

func Getselector() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}){
	return FibonacciinvokerM{}.Selector
}

func NewFibonacciInvokerM() FibonacciinvokerM {
	r := new(FibonacciinvokerM)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (e FibonacciinvokerM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	e.I_Process(msg, info)
}

func (FibonacciinvokerM) I_Process(msg *messages.SAMessage, info [] *interface{}) { // TODO

	// unmarshall
	payload := msg.Payload.([]byte)

	miopPacket := miop.Packet{}
	err := msgpack.Unmarshal(payload, &miopPacket)
	if err != nil {
		log.Fatalf("NamingInvokerM:: %s", err)
	}

	//var inv shared.Request
	op := miopPacket.Bd.ReqHeader.Operation
	switch op {
	case "Fibo":
		_p0 := int(miopPacket.Bd.ReqBody.Body[0].(int64))

		fmt.Printf("FibonacciInvokerM [Plugin V2]\n")

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
			log.Fatalf("NamingInvokerM:: %s", err)
		}

		toSRH := make([]interface{}, 1, 1)
		toSRH[0] = r

		*msg = messages.SAMessage{Payload: toSRH}
	default:
		fmt.Printf("FibonacciInvokerM:: Operation '%v' not implemented by Fibonacci Service\n", op)
		os.Exit(0)
	}
}
