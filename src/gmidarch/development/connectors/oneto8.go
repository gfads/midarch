package connectors

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type Oneto8 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto8() Oneto8 {

	r := new(Oneto8)
	//r.Behaviour = "B = InvP.e1 -> InvR.e2 -> I_Debug1 -> InvR.e3 -> InvR.e4 -> InvR.e5 -> InvR.e6 -> InvR.e7 -> InvR.e8 -> InvR.e9 -> I_Debug2 -> B"
	r.Behaviour = "B = InvP.e1 -> InvR.e2 -> P2 [] I_Timeout -> P2 \n P2 = InvR.e3 -> P3 [] I_Timeout -> P3 \n P3 = InvR.e4 -> P4 [] I_Timeout -> P4 \n P4 = InvR.e5 -> P5 [] I_Timeout -> P5 \n P5 = InvR.e6 -> P6 [] I_Timeout -> P6 \n P6 = InvR.e7 -> P7 [] I_Timeout -> P7 \n P7 = InvR.e8 -> P8 [] I_Timeout -> P8 \n P8 = InvR.e9 -> B [] I_Timeout -> B"
	//r.Behaviour = "B = InvP.e1 -> I_Debug1 -> InvR.e2 -> P2 [] I_Timeout -> P2 \n P2 = InvR.e3 -> P3 [] I_Timeout -> P3 \n P3 = InvR.e4 -> P4 [] I_Timeout -> P4 \n P4 = InvR.e5 -> P5 [] I_Timeout -> P5 \n P5 = InvR.e6 -> P6 [] I_Timeout -> P6 \n P6 = InvR.e7 -> P7 [] I_Timeout -> P7 \n P7 = InvR.e8 -> P8 [] I_Timeout -> P8 \n P8 = InvR.e9 -> I_Debug2 -> Oneto8 [] I_Timeout -> Oneto8\n"

	return *r
}

func (e Oneto8) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Timeout(msg, info)
}

func (Oneto8) I_Timeout(msg *messages.SAMessage, info [] *interface{}) {
	*msg = *msg
}

func (Oneto8) I_Debug1(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Oneto8:: I_Debug1 :: %v\n", msg.Payload)
}

func (Oneto8) I_Debug2(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Printf("Oneto8:: I_Debug2 :: %v\n", msg.Payload)
}
