package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Oneto9 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto9() Oneto9 {

	r := new(Oneto9)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> InvR.e3 -> InvR.e4 -> InvR.e5 -> InvR.e6 -> InvR.e7 -> InvR.e8 -> InvR.e9 -> InvR.e10 -> B"

	return *r
}