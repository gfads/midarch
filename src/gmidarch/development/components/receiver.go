package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type Receiver struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewReceiver() Receiver {
	r := new(Receiver)
	r.Behaviour = "B = InvP.e1 -> I_PrintMessage -> B"

	return *r
}

func (Receiver) Selector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}){

	var f func(*messages.SAMessage,[]*interface{})
	switch op {
	case "I_Printmessage":
		f = func(msg *messages.SAMessage, info []*interface{}) {
			elem.(Receiver).I_Printmessage(msg, info)
		}
	}
	return f
}

func (Receiver) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Receiver:: %v  \n", *msg)
}
