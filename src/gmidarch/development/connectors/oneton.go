package connectors

import (
	graphs2 "gmidarch/development/artefacts/graphs"
	"shared/shared"
)

type OnetoN struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewOnetoN() OnetoN {

	// create a new instance of client
	r := new(OnetoN)
	r.Behaviour = "B = "+ shared.RUNTIME_BEHAVIOUR

	return *r
}