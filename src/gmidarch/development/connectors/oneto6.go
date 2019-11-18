package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Oneto6 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto6() Oneto6 {

	r := new(Oneto6)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> InvR.e3 -> InvR.e4 -> InvR.e5 -> InvR.e6 -> InvR.e7 -> B"

	return *r
}