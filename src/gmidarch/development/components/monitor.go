package components

import (
	graphs2 "gmidarch/development/artefacts/graphs"
	messages2 "gmidarch/development/messages"
)

type Monitor struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewMonitor() Monitor {

	// create a new instance of Server
	r := new(Monitor)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (Monitor) I_Process(msg *messages2.SAMessage, info [] *interface{}) {
	*msg = messages2.SAMessage{Payload: msg.Payload}
}
