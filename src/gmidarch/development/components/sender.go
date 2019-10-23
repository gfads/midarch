package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type Sender struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewSender() Sender {

	r := new(Sender)
	r.Behaviour = "B = I_Setmessage1 -> InvR.e1 -> B [] I_Setmessage2 -> InvR.e1 -> B [] I_Setmessage3 -> InvR.e1 -> B"

	return *r

}

func (Sender) I_Setmessage1(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 1)"}
}
func (Sender) I_Setmessage2(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 2)"}
}
func (Sender) I_Setmessage3(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 3)"}
}
func (Sender) I_Debug(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Sender:: Debug:: %v \n", *msg)
}
