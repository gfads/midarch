package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
	"time"
)

var idx int
var t1, t2 time.Time
var durations [shared.SAMPLE_SIZE] time.Duration

type Fibonacciclient struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewFibonacciclient() Fibonacciclient {

	r := new(Fibonacciclient)
	r.Behaviour = "B = I_Setmessage -> InvR.e1 -> TerR.e1 -> I_Printmessage -> B"

	return *r
}

func (e Fibonacciclient) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'S' { // Set message
		e.I_Setmessage(msg, info)
	} else { // "I_Printmessage":
		e.I_Printmessage(msg, info)
	}
}

func (Fibonacciclient) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {
	argsTemp := make([]interface{}, 1, 1)

	t1 = time.Now() // start time
	argsTemp[0] = 10 // TODO
	*msg = messages.SAMessage{Payload: shared.Request{Op: "fibo", Args: argsTemp}}
}

func (Fibonacciclient) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
    fmt.Printf("Fibonacciclient:: %v\n",msg.Payload)
}
