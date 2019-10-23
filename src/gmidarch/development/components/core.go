package components

import (
	"fmt"
	graphs2 "gmidarch/development/artefacts/graphs"
	messages2 "gmidarch/development/messages"
	"os"
	shared2 "shared"
)

var count int

type Core struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewCore() Core {

	// create a new instance of Server
	r := new(Core)
	r.Behaviour = "B = "+ shared2.RUNTIME_BEHAVIOUR
	//r.Behaviour = "B = InvR.e1 -> P1\n P1 = InvP.e2 -> P1"

	return *r
}

func (Core) I_Debug(msg *messages2.SAMessage, info [] *interface{}) {
	fmt.Printf("******************* Core:: I_Debug ****************** \n")

	count++

	if count == 5 {
		os.Exit(0)
	}

}
