package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

type Fibonacciproxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewFibonacciproxy() Fibonacciproxy {

	r := new(Fibonacciproxy)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e Fibonacciproxy) Selector(elem interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { //"I_Out"
		e.I_Out(msg, info)
	}
}

func (Fibonacciproxy) I_In(msg *messages.SAMessage, info [] *interface{}) {
	inv := shared.Invocation{Host: "localhost", Port: shared.FIBONACCI_PORT, Req: msg.Payload.(shared.Request)} // TODO

	*msg = messages.SAMessage{Payload: inv}
}

func (Fibonacciproxy) I_Out(msg *messages.SAMessage, info [] *interface{}) {

	result := msg.Payload.([]interface{})
	*msg = messages.SAMessage{Payload: result[0]}
}
