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
	r.Behaviour = "B = InvP.e1 -> I_ProcessIn -> InvR.e2 -> TerR.e2 -> I_ProcessOut -> TerP.e1 -> B"

	return *r
}

func (Fibonacciproxy) Selector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}) {

	var f func(*messages.SAMessage, []*interface{})
	if op == "I_Processin" {
		f = func(msg *messages.SAMessage, info []*interface{}) {
			elem.(Fibonacciproxy).I_Processin(msg, info)
		}
	} else { //"I_Processout"
			f = func(msg *messages.SAMessage, info []*interface{}) {
				elem.(Fibonacciproxy).I_Processout(msg, info)
			}
		}
	return f
}

func (Fibonacciproxy) I_Processin(msg *messages.SAMessage, info [] *interface{}) {
	inv := shared.Invocation{Host:"localhost",Port:shared.FIBONACCI_PORT,Req:msg.Payload.(shared.Request)} // TODO

	*msg = messages.SAMessage{Payload: inv}
}

func (Fibonacciproxy) I_Processout(msg *messages.SAMessage, info [] *interface{}) {

	result := msg.Payload.([]interface{})
	*msg = messages.SAMessage{Payload: result[0]}
}
