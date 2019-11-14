package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Oneto2 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto2() Oneto2 {

	r := new(Oneto2)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> InvR.e3 -> B"

	return *r
}