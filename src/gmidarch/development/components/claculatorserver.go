package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"shared"
)

type Calculatorserver struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Newcalculatorserver() Calculatorserver {

	// create a new instance of Server
	r := new(Calculatorserver)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (Calculatorserver) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared.Request)
	op := req.Op
	p1 := req.Args[0].(int)
	p2 := req.Args[1].(int)
	r := 0

	switch op {
	case "add":
		r = p1+p2
	case "sub":
		r = p1-p2
	case "mul":
		r = p1 * p2
	case "div":
		r = p1/p2
	default:
		fmt.Println("Server Calculator:: Operation '" + op + "' not supported!!")
		os.Exit(0)
	}

	*msg = messages.SAMessage{Payload: r}
}