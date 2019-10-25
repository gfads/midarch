package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"shared"
)

var count int

type Core struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewCore() Core {

	r := new(Core)
	r.Behaviour = "B = "+ shared.RUNTIME_BEHAVIOUR

	return *r
}

func (Core) I_Debug(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("******************* Core:: I_Debug ****************** \n")

	count++

	if count == 5 {
		os.Exit(0)
	}

}
