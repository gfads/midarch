package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Oneto8 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto8() Oneto8 {

	r := new(Oneto8)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> InvR.e3 -> InvR.e4 -> InvR.e5 -> InvR.e6 -> InvR.e7 -> InvR.e8 -> InvR.e9 -> B"

	return *r
}