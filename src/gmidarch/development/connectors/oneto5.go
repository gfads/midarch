package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Oneto5 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto5() Oneto5 {

	r := new(Oneto5)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> InvR.e3 -> InvR.e4 -> InvR.e5 -> InvR.e6 -> B"

	return *r
}