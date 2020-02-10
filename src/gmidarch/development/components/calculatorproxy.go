package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"shared"
)

type CalculatorProxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewCalculatorProxy() CalculatorProxy {

	// create a new instance of Server
	r := new(CalculatorProxy)
	r.Behaviour = "B = InvP.e1 -> I_ProcessIn -> InvR.e2 -> TerR.e2 -> I_ProcessOut -> TerP.e1 -> B"

	return *r
}

func (CalculatorProxy) I_Processin(msg *messages.SAMessage, info [] *interface{}) {
	inv := shared.Invocation{}
	inv.Host = "localhost"             // TODO
	inv.Port = shared.CALCULATOR_PORT // TODO
	inv.Req = msg.Payload.(shared.Request)

	*msg = messages.SAMessage{Payload: inv}
}

func (CalculatorProxy) I_Processout(msg *messages.SAMessage, info [] *interface{}) {

	result := msg.Payload.([]interface{})
	*msg = messages.SAMessage{Payload: int(result[0].(float64))}
}
