package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
)

type Analyser struct {
	Behaviour   string
	Graph graphs.ExecGraph
}

func NewAnalyser() Analyser {

	// create a new instance of Server
	r := new(Analyser)
	r.Behaviour = "B = InvP.e1 -> I_Process -> InvR.e2 -> B"

	return *r
}

func (Analyser) I_Process(msg *messages.SAMessage,info [] *interface{}) {
	//fmt.Printf("Analyser:: I_Process \n")
	*msg = messages.SAMessage{Payload:"TODO"} //TODO
}
