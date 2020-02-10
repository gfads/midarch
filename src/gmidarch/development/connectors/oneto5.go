package connectors

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type Oneto5 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto5() Oneto5 {

	r := new(Oneto5)
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> P2 [] I_Timeout -> P2 \n P2 = InvR.e3 -> P3 [] I_Timeout -> P3 \n P3 = InvR.e4 -> P4 [] I_Timeout -> P4 \n P4 = InvR.e5 -> P5 [] I_Timeout -> P5 \n P5 = InvR.e6 -> B [] I_Timeout -> B"

	return *r
}

func (e Oneto5) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Timeout(msg, info)
}

func (Oneto5) I_Timeout(msg *messages.SAMessage, info [] *interface{}) {
	*msg = *msg
}

func (Oneto5) I_Debug1(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Oneto5:: I_Debug1 :: %v\n", msg.Payload)
}

func (Oneto5) I_Debug2(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Oneto5:: I_Debug2 :: %v\n", msg.Payload)
}
