package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
)

type Executor struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewExecutor() Executor {

	// create a new instance of client
	r := new(Executor)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (Executor) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload:"TODO"}
	//fmt.Printf("Executor:: I_Process \n")
}
