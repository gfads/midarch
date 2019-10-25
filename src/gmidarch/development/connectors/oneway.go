package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Oneway struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneway() Oneway {

	// create a new instance of client
	r := new(Oneway)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> B"

	return *r
}