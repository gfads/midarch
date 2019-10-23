package main

import (
	"fmt"
	graphs2 "gmidarch/development/artefacts/graphs"
	messages2 "gmidarch/development/messages"
)

type Receiver struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func Gettype() interface{} {
	return Receiver{}
}

func NewReceiver() Receiver {

	// create a new instance of client
	r := new(Receiver)
	r.Behaviour = "B = InvP.e1 -> I_PrintMessage -> B"

	return *r
}

func (Receiver) I_Printmessage(msg *messages2.SAMessage, info [] *interface{}) {
	fmt.Printf("Receiver:: Plugin [V1]:: %v  \n", *msg)
}