package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
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

func (Client) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World from Client"}
}

func (Client) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {

	fmt.Println(msg.Payload)
}