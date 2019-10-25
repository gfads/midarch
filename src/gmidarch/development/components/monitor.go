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

func (Monitor) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: msg.Payload}
}
