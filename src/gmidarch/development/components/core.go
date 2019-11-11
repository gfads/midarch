package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

var count int

type Core struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewCore() Core {

	r := new(Core)
	r.Behaviour = "B = " + shared.RUNTIME_BEHAVIOUR

	return *r
}

func (Core) Selector(elem interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
}
