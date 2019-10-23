package main

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
)

type Receiver struct {
	Behaviour string
	Graph     graphs.ExecGraph
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

func (Receiver) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Receiver:: Plugin:: %v  \n", *msg)
}