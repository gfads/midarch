package middleware

import (
	"apps/businesses/fibonacciImpl"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	"shared"
)

//@Type: FibonacciInvoker
//@Behaviour: Behaviour = InvP.e1 -> I_Beforeunmarshalling -> InvR.e2 -> TerR.e2 -> I_Beforeserver -> I_Beforemarshalling -> InvR.e2 -> TerR.e2 -> I_Beforesend -> TerP.e1 -> Behaviour
type FibonacciInvoker struct{}

func (FibonacciInvoker) I_Beforeunmarshalling(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	tempParams := []interface{}{msg.Payload}
	msg.Payload = messages.FunctionalRequest{Op: "unmarshall", Params: tempParams}
}

func (FibonacciInvoker) I_Beforeserver(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	miopPacket := msg.Payload.(messages.FunctionalReply).Rep.(miop.MiopPacket) // from marshaller

	req := messages.FunctionalRequest{Op:miopPacket.Bd.ReqHeader.Operation, Params:miopPacket.Bd.ReqBody.Body}

	switch req.Op {
	case "F":

		// Parameters
		params := []interface{}{req.Params[0].(int)}

		// Functional request
		req2 := messages.FunctionalRequest{Op: req.Op, Params: params}
		msg.Payload = &req2

		reply := fibonacciImpl.Fibonacci{}.F(req.Params[0].(int))
		msg.Payload = messages.FunctionalReply{Rep: reply}

	default:
		shared.ErrorHandler(shared.GetFunction(), "Operation '"+req.Op+"' not present in Invoker")
	}
}

func (FibonacciInvoker) I_Beforemarshalling(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	reply := msg.Payload.(messages.FunctionalReply)
	repPacket := miop.CreateRepPacket(reply.Rep)
	tempParams := []interface{}{repPacket}
	msg.Payload = messages.FunctionalRequest{Op: "marshall", Params: tempParams}
}

func (FibonacciInvoker) I_Beforesend(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	msg.Payload = msg.Payload.(messages.FunctionalReply).Rep
}
