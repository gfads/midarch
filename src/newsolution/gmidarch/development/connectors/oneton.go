package connectors

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/shared/parameters"
)

type OnetoN struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOnetoN() OnetoN {

	// create a new instance of client
	r := new(OnetoN)
	r.Behaviour = "B = "+parameters.RUNTIME_BEHAVIOUR

	return *r
}