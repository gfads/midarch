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

func (e Fibonacciclient) Selector(elem interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	if op[2] == 'S' { // Set message
		e.I_Setmessage(msg, info)
	} else { // "I_Printmessage":
		e.I_Printmessage(msg, info)
	}
}

func (Fibonacciclient) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {
	argsTemp := make([]interface{}, 1, 1)

	t1 = time.Now() // start time
	argsTemp[0] = 1
	*msg = messages.SAMessage{Payload: shared.Request{Op: "fibo", Args: argsTemp}}
}

func (Fibonacciclient) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	//fmt.Printf("Fibonacciclient:: %v [%v]\n",msg.Payload,idx)

	t2 = time.Now() // finish time

//	if idx > 1500 {
//		fmt.Println(float64(t2.Sub(t1).Nanoseconds()) / float64(1000000))
//	}

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

	//fmt.Printf("%v\n",durations[idx])
	idx++
}
