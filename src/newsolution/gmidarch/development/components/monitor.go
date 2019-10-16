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
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (Monitor) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	//fmt.Printf("Monitor:: I_Process\n")
	*msg = messages.SAMessage{Payload:"TODO"}
}
