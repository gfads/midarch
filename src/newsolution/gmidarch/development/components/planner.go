package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
)

type Planner struct {
	Behaviour   string
	Graph graphs.ExecGraph
}

func NewPlanner() Planner {

	// create a new instance of Server
	r := new(Planner)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (Planner) I_Process(msg *messages.SAMessage,info [] *interface{}) {
	//fmt.Printf("Planner:: I_Process \n")
	*msg = messages.SAMessage{} // TODO
}
