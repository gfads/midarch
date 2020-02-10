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
	//r.Behaviour = "B = I_Setmessage1 -> InvR.e1 -> B"
	return *r
}

func (s Sender) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {

	switch op {
	case "I_Setmessage1":
		elem.(Sender).I_Setmessage1(msg, info)
	case "I_Setmessage2":
		elem.(Sender).I_Setmessage2(msg, info)
	case "I_Setmessage3":
		elem.(Sender).I_Setmessage3(msg, info)
	}
}

func (s Sender) OldSelector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}) {

	var f func(*messages.SAMessage, []*interface{})
	switch op {
	case "I_Setmessage1":
		f = func(msg *messages.SAMessage, info []*interface{}) {
			elem.(Sender).I_Setmessage1(msg, info)
		}
	case "I_Setmessage2":
		f = func(msg *messages.SAMessage, info []*interface{}) {
			elem.(Sender).I_Setmessage2(msg, info)
		}
	case "I_Setmessage3":
		f = func(msg *messages.SAMessage, info []*interface{}) {
			elem.(Sender).I_Setmessage3(msg, info)
		}
	}
	return f
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
