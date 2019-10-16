package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
)

type Monevolutive struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewMonevolutive() Monevolutive {

	// create a new instance of Server
	r := new(Monevolutive)
	r.Behaviour = "B = I_Collect -> InvR.e1 -> B"

	return *r
}

func (Monevolutive) I_Collect(msg *messages.SAMessage, info [] *interface{}) {
	//fmt.Printf("MonitorEvolutive:: I_Collect \n")
	*msg = messages.SAMessage{"TODO"} // TODO
}
