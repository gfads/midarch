package main

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

type Calculatorclient struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Gettype() interface{} {
	return Calculatorclient{}
}

func NewCalculatorclient() Calculatorclient {

	r := new(Calculatorclient)
	r.Behaviour = "B = I_Setmessage -> InvR.e1 -> TerR.e1 -> I_Printmessage -> B"

	return *r
}

func (Calculatorclient) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {

	argsTemp := make([]interface{}, 2)
	argsTemp[0] = 1
	argsTemp[1] = 2
	*msg = messages.SAMessage{Payload: shared.Request{Op: "add", Args: argsTemp}}

}

func (Calculatorclient) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Println(msg.Payload)
}
