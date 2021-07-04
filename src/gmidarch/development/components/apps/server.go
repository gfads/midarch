package apps

import (
	"apps/businesses/calculatorimpl"
	"gmidarch/development/messages"
	"shared"
)

//@Type: Server
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type Server struct{}

func (s Server) I_Process(id string, msg *messages.SAMessage, info *interface{}) {
	request := msg.Payload.(*messages.FunctionalRequest)

	switch request.Op {
	case "Add":
		reply := calculatorimpl.CalculatorImpl{}.Add(request.Params[0].(int), request.Params[1].(int))
		msg.Payload = messages.FunctionalReply{Rep:reply}
	default:
		shared.ErrorHandler(shared.GetFunction(), "Operation '"+request.Op+ "' not implemented by Calculator")
	}
}
