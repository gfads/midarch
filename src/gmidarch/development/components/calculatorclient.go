package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"shared"
	"time"
)

type Calculatorclient struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewCalculatorclient() Calculatorclient {

	r := new(Calculatorclient)
	r.Behaviour = "B = I_Setmessage -> InvR.e1 -> TerR.e1 -> I_Printmessage -> B"

	return *r
}

func (Calculatorclient) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {

	t1 = time.Now()  // start time

	argsTemp := make([]interface{}, 2)
	argsTemp[0] = 1
	argsTemp[1] = 2
	*msg = messages.SAMessage{Payload: shared.Request{Op: "add", Args: argsTemp}}
}

func (Calculatorclient) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	//fmt.Printf("Calculatorclient:: %v [%v]\n",msg.Payload,idx)

	t2 := time.Now()   // finish time

	if idx >= shared.SAMPLE_SIZE {
		totalTime := time.Duration(0)
		for i := range durations {
			totalTime += durations[i]
		}
		fmt.Printf("Total Time [%v]: %v \n", shared.SAMPLE_SIZE, totalTime)
		os.Exit(0)
	}
	durations[idx] = t2.Sub(t1)
	idx++
}
