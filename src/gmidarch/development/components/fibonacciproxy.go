package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

type Fibonacciproxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      string
}

var i_PreInvRFP = make(chan messages.SAMessage)
var i_PosTerRFP = make(chan messages.SAMessage)

func NewFibonacciproxy() Fibonacciproxy {

	r := new(Fibonacciproxy)
	r.Behaviour = "B = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> B"

	return *r
}

func (e Fibonacciproxy) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { //"I_Out"
		e.I_Out(msg, info)
	}
}

func (e Fibonacciproxy) Fibo(_p1 int) int {
	_args := []interface{}{_p1}
	_reqMsg := messages.SAMessage{messages.Invocation{Host: e.Host, Port: e.Port, Op: "Fibo", Args: _args}}

	i_PreInvRFP  <- _reqMsg
	_repMsg := <-i_PosTerRFP

	_reply := int(_repMsg.Payload.(int64))

	return _reply
}

func (Fibonacciproxy) I_In(msg *messages.SAMessage, info [] *interface{}) {
	*msg = <- i_PreInvRFP
}

func (Fibonacciproxy) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	i_PosTerRFP <- *msg
}

func (Fibonacciproxy) I_InOld(msg *messages.SAMessage, info [] *interface{}) {
	inv := shared.Invocation{Host: "localhost", Port: shared.FIBONACCI_PORT, Req: msg.Payload.(shared.Request)} // TODO

	*msg = messages.SAMessage{Payload: inv}
}

func (Fibonacciproxy) I_OutOld(msg *messages.SAMessage, info [] *interface{}) {

	result := msg.Payload.([]interface{})
	*msg = messages.SAMessage{Payload: result[0]}
}
