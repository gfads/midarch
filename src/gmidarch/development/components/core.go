package components

import (
	"fmt"
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

func (Core) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
//	Core{}.I_Debug(msg,info)
}

func (Core) I_Debug(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Core:: %v\n",msg.Payload)
}
