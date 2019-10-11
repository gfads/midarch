package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
)

type Core struct {
	Behaviour   string
	Graph graphs.ExecGraph
}

func NewCore() Core {

	// create a new instance of Server
	r := new(Core)
	//r.Behaviour = "B = "+parameters.RUNTIME_BEHAVIOUR
	r.Behaviour = "B = InvR.e1 -> B"

	return *r
}