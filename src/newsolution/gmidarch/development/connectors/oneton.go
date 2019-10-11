package connectors

import "newsolution/gmidarch/development/artefacts/graphs"

type OnetoN struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOnetoN() OnetoN {

	// create a new instance of client
	r := new(OnetoN)
	r.Behaviour = "Oneway = InvP.e1 -> InvR.e2 -> Oneway"

	return *r
}