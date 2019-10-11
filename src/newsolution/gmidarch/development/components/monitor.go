package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
)

type Monitor struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewMonitor() Monitor {

	// create a new instance of Server
	r := new(Monitor)
	r.Behaviour = "B = I_Collect -> InvR.e1 -> B"

	return *r
}

func (Monitor) I_Collect(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{} // TODO
}
