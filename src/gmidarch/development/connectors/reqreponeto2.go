package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Reqrepeoneto2 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Newreqreponeto2() Reqrepeoneto2 {

	// create a new instance of client
	r := new(Reqrepeoneto2)
	r.Behaviour = "B = InvP.e1 -> (InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] InvR.e3 -> TerR.e3 -> TerP.e1 -> B)"

	return *r
}