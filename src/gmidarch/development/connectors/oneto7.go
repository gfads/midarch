package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Oneto7 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto7() Oneto7 {

	r := new(Oneto7)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> InvR.e3 -> InvR.e4 -> InvR.e5 -> InvR.e6 -> InvR.e7 -> InvR.e8 -> B"

	return *r
}