package components

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/parameters"
	"os"
)

var count int

type Core struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewCore() Core {

	// create a new instance of Server
	r := new(Core)
	r.Behaviour = "B = "+parameters.RUNTIME_BEHAVIOUR
	//r.Behaviour = "B = InvR.e1 -> P1\n P1 = InvP.e2 -> P1"

	return *r
}

func (Core) I_Debug(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("******************* Core:: I_Debug ****************** \n")

	count++

	if count == 5 {
		os.Exit(0)
	}

}
