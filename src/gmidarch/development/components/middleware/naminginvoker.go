package middleware

import (
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	"reflect"
	"shared"
)

//@Type: Naminginvoker
//@Behaviour: Behaviour = InvP.e1 -> I_Beforeunmarshalling -> InvR.e2 -> TerR.e2 -> I_Beforeserver -> InvR.e3 -> TerR.e3 -> I_Beforemarshalling -> InvR.e2 -> TerR.e2 -> I_Beforesend -> TerP.e1 -> Behaviour
type Naminginvoker struct{}

func (Naminginvoker) I_Beforeunmarshalling(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	tempParams := []interface{}{msg.Payload}
	msg.Payload = messages.FunctionalRequest{Op: "unmarshall", Params: tempParams}
}

func (Naminginvoker) I_Beforeserver(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	miopPacket := msg.Payload.(messages.FunctionalReply).Rep.(miop.MiopPacket) // from marshaller

	op := miopPacket.Bd.ReqHeader.Operation
	p := miopPacket.Bd.ReqBody.Body

	req := messages.FunctionalRequest{Op: op, Params: p}

	switch req.Op {
	case "Register":
		if reflect.TypeOf(req.Params[1]).String() == "map[string]interface {}" { // JSON
			// AOR
			aorMap := req.Params[1].(map[string]interface{})
			aorHost := aorMap["host"].(string)
			aorPort := aorMap["port"].(string)
			aorId := int(aorMap["id"].(float64))
			aorProxy := aorMap["proxy"].(string)
			aor := messages.AOR{Host: aorHost, Port: aorPort, Id: aorId, ProxyName: aorProxy}

			// Parameters
			params := []interface{}{}
			params = append(params, req.Params[0].(string))
			params = append(params, aor)

			// Functional request
			req1 := messages.FunctionalRequest{Op: req.Op, Params: params}
			msg.Payload = &req1

		} else { // GOB
			// Parameters
			params := []interface{}{}
			params = append(params, req.Params[0].(string))
			params = append(params, req.Params[1].(messages.AOR))

			// Functional request
			req1 := messages.FunctionalRequest{Op: req.Op, Params: params}
			msg.Payload = &req1
		}

	case "Lookup":

		// Parameters
		params := []interface{}{req.Params[0].(string), req.Params[0].(string)}

		// Functional request
		req1 := messages.FunctionalRequest{Op: req.Op, Params: params}
		msg.Payload = &req1

	case "List":

		// Parameters
		params := []interface{}{}

		// Functional request
		req1 := messages.FunctionalRequest{Op: req.Op, Params: params}
		msg.Payload = &req1

	default:
		if req.Op != "" {
			shared.ErrorHandler(shared.GetFunction(), "Operation '"+req.Op+"' not present in Naming Invoker")
		}
	}
}

func (Naminginvoker) I_Beforemarshalling(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	reply := msg.Payload.(messages.FunctionalReply)
	repPacket := miop.CreateRepPacket(reply.Rep)
	tempParams := []interface{}{repPacket}
	msg.Payload = messages.FunctionalRequest{Op: "marshall", Params: tempParams}
}

func (Naminginvoker) I_Beforesend(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	msg.Payload = msg.Payload.(messages.FunctionalReply).Rep
}
