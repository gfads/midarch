package main

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type Receiver struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Gettype() interface{} {
	return Receiver{}
}

func Getselector() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}){
	return Receiver{}.Selector
}

func NewReceiver() Receiver {

	// create a new instance of client
	r := new(Receiver)
	r.Behaviour = "B = InvP.e1 -> I_PrintMessage -> B"

	return *r
}

func (Receiver) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}){
	elem.(Receiver).I_Printmessage(msg, info)
}

func (Receiver) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Receiver:: Plugin [V1]:: %v  \n", *msg)
}