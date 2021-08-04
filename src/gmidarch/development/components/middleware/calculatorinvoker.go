package middleware

import (
	"apps/businesses/calculatorimpl"
	"fmt"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	"shared"
)

//@Type: Calculatorinvoker
//@Behaviour: Behaviour = InvP.e1 -> I_Beforeunmarshalling -> InvR.e2 -> TerR.e2 -> I_Beforeserver -> I_Beforemarshalling -> InvR.e2 -> TerR.e2 -> I_Beforesend -> TerP.e1 -> Behaviour
type Calculatorinvoker struct{}

func (Calculatorinvoker) I_Beforeunmarshalling(id string, msg *messages.SAMessage, info *interface{}) {
	tempParams := []interface{}{msg.Payload}
	msg.Payload = messages.FunctionalRequest{Op: "unmarshall", Params: tempParams}
}

func (Calculatorinvoker) I_Beforeserver(id string, msg *messages.SAMessage, info *interface{}) {
	miopPacket := msg.Payload.(messages.FunctionalReply).Rep.(miop.MiopPacket) // from marshaller

	fmt.Println(shared.GetFunction(),"TO AQUI")
	req := messages.FunctionalRequest{Op:miopPacket.Bd.ReqHeader.Operation, Params:miopPacket.Bd.ReqBody.Body}

	switch req.Op {
	case "Add":

		// Parameters
		params := []interface{}{req.Params[0].(int), req.Params[1].(int)}

		// Functional request
		req2 := messages.FunctionalRequest{Op: req.Op, Params: params}
		msg.Payload = &req2

		reply := calculatorimpl.CalculatorImpl{}.Add(req.Params[0].(int), req.Params[1].(int))
		msg.Payload = messages.FunctionalReply{Rep:reply}

	default:
		shared.ErrorHandler(shared.GetFunction(), "Operation '"+req.Op+"' not presente in Invoker")
	}
}

func (Calculatorinvoker) I_Beforemarshalling(id string, msg *messages.SAMessage, info *interface{}) {
	reply := msg.Payload.(messages.FunctionalReply)
	repPacket := miop.CreateRepPacket(reply.Rep)
	tempParams := []interface{}{repPacket}
	msg.Payload = messages.FunctionalRequest{Op: "marshall", Params: tempParams}
}

func (Calculatorinvoker) I_Beforesend(id string, msg *messages.SAMessage, info *interface{}) {
	msg.Payload = msg.Payload.(messages.FunctionalReply).Rep
}
