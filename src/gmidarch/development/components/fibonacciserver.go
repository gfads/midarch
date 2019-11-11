package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

type Fibonacciserver struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Newfibonacciserver() Fibonacciserver {

	r := new(Fibonacciserver)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (e Fibonacciserver) Selector(elem interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	e.I_Process(msg, info)
}

func (Fibonacciserver) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared.Request)

		*msg = messages.SAMessage{Payload: f(req.Args[0].(int))}
}

func f(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return f(n-1) + f(n-2)
	}
}
