package connectors

import (
	"gmidarch/development/artefacts/graphs"
)

type Requestreply struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewRequestReply() Requestreply {

	// create a new instance of client
	r := new(Requestreply)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B"

	return *r
}