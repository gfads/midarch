package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"shared"
	"time"
)

type Client struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewClient() Client {

	r := new(Client)
	r.Behaviour = "B = I_Setmessage -> InvR.e1 -> TerR.e1 -> I_Printmessage -> B"

	return *r
}

func (Client) Selector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}){

	var f func(*messages.SAMessage,[]*interface{})
	switch op {
	case "I_Setmessage":
		f = func(msg *messages.SAMessage, info []*interface{}){
			elem.(Client).I_Setmessage(msg,info)
		}
	case "I_Printmessage":
		f = func(msg *messages.SAMessage, info []*interface{}){
			elem.(Client).I_Printmessage(msg,info)
		}
	}
	return f
}

var n int = 0
var duration [2500] time.Duration

func (Client) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {
	t1 = time.Now()
	*msg = messages.SAMessage{Payload: "Hello World from Client"}
}

func (Client) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Println(msg.Payload)

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