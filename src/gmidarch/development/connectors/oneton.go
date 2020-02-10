package connectors

import (
	"gmidarch/development/artefacts/graphs"
	"shared"
)

type OnetoN struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOnetoN() OnetoN {

	// create a new instance of client
	r := new(OnetoN)
	r.Behaviour = "B = "+ shared.RUNTIME_BEHAVIOUR

	return *r
}