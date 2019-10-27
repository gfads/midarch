package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"shared"
	"time"
)

var times [1000] time.Duration
var idx int
var t1 time.Time

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

	time.Sleep(100 * time.Millisecond)

	//if idx < 100 {
	t1 = time.Now()
	argsTemp := make([]interface{}, 2)
	argsTemp[0] = 1
	argsTemp[1] = 2
	*msg = messages.SAMessage{Payload: shared.Request{Op: "add", Args: argsTemp}}
	//}
}

func (Calculatorclient) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Calculatorclient:: %v [%v]\n",msg.Payload,idx)
	//times[idx] = time.Now().Sub(t1)

	if idx >= 10000 {
		fmt.Printf("CalculatorClient:: Experiment finished!")
		os.Exit(0)
	}
	idx++
}
