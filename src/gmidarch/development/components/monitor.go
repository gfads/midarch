package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type Monitor struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewMonitor() Monitor {

	// create a new instance of Server
	r := new(Monitor)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (e Monitor) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Process(msg, info)
}

func (Monitor) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: msg.Payload}
}
