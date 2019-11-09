package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"shared"
	"time"
)

var idx int
var t1 time.Time
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

func (Fibonacciclient) Selector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}) {

	if op == "I_Setmessage" {
		return func(msg *messages.SAMessage, info []*interface{}) {
			elem.(Fibonacciclient).I_Setmessage(msg, info)
		}
	} else { // "I_Printmessage":
		return func(msg *messages.SAMessage, info []*interface{}) {
			elem.(Fibonacciclient).I_Printmessage(msg, info)
		}
	}
}

func (Fibonacciclient) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {

	t1 = time.Now() // start time
	argsTemp := make([]interface{}, 1)
	argsTemp[0] = 0
	*msg = messages.SAMessage{Payload: shared.Request{Op: "fibo", Args: argsTemp}}
}

func (Fibonacciclient) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	//fmt.Printf("Fibonacciclient:: %v [%v]\n",msg.Payload,idx)

	t2 := time.Now() // finish time

	if idx == shared.SAMPLE_SIZE {
		totalTime := time.Duration(0)
		for i := range durations {
			totalTime += durations[i]
		}

		fmt.Printf("Total   Time [%v]: %v \n", shared.SAMPLE_SIZE, totalTime)
		fmt.Printf("Average Time [%v]: %v \n", shared.SAMPLE_SIZE, totalTime/shared.SAMPLE_SIZE)

		os.Exit(0)
	}
	durations[idx] = t2.Sub(t1)
	idx++
}
