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

func (Core) Selector(elem interface{}, op string, msg *messages.SAMessage, info []*interface{}){


			elem.(Core).I_Debug(msg,info)
}

func (Core) OldSelector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}){

	var f func(*messages.SAMessage,[]*interface{})
	switch op {
	case "I_Debug":
		f = func(msg *messages.SAMessage, info []*interface{}){
			elem.(Core).I_Debug(msg,info)
		}
	}
	return f
}

func (Core) I_Debug(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("******************* Core:: I_Debug ****************** \n")

	count++

	if count == 5 {
		os.Exit(0)
	}

}
