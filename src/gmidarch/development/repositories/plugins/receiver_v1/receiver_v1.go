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

func Getselector() func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool){
	return Receiver{}.Selector
}

func NewReceiver() Receiver {

	// create a new instance of client
	r := new(Receiver)
	r.Behaviour = "B = InvP.e1 -> I_PrintMessage -> B"

	return *r
}

func (e Receiver) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool){
	e.I_Printmessage(msg,r)
}

func (Receiver) I_Printmessage(msg *messages.SAMessage, r *bool) {
	fmt.Printf("Receiver:: Plugin [V1]:: %v  \n", *msg)
}