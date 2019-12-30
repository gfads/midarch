package connectors

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
)

type Oneto9 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewOneto9() Oneto9 {

	r := new(Oneto9)
	r.Behaviour = "B = InvP.e1 -> (InvR.e2 -> P2 [] I_Timeout -> P2) \n P2 = InvR.e3 -> P3 [] I_Timeout -> P3 \n P3 = InvR.e4 -> P4 [] I_Timeout -> P4 \n P4 = InvR.e5 -> P5 [] I_Timeout -> P5 \n P5 = InvR.e6 -> P6 [] I_Timeout -> P6 \n P6 = InvR.e7 -> P7 [] I_Timeout -> P7 \n P7 = InvR.e8 -> P8 [] I_Timeout -> P8 \n P8 = InvR.e9 -> P9 [] I_Timeout -> P9 \n P9 = InvR.e10 -> B [] I_Timeout -> B"

	return *r
}

func (e Oneto9) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	e.I_Timeout(msg, info)
}

func (Oneto9) I_Timeout(msg *messages.SAMessage, info [] *interface{}) {
	*msg = *msg
}
