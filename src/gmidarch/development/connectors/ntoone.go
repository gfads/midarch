package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type NtoOne struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewNtoOne() NtoOne {

	// create a new instance of client
	r := new(NtoOne)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] InvP.e3 -> InvR.e2 -> TerR.e2 -> TerP.e3 -> B"

	return *r
}