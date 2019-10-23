package components

import (
	graphs2 "gmidarch/development/artefacts/graphs"
)

type ComponentNovo struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewComponentNovo() ComponentNovo {

	// create a new instance of Server
	r := new(ComponentNovo)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}